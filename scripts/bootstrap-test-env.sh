#!/usr/bin/env bash
# Copyright 2025 InferaDB
# SPDX-License-Identifier: Apache-2.0
#
# bootstrap-test-env.sh - Bootstrap test environment for Terraform Provider acceptance tests
#
# This script:
# 1. Waits for the Control service to be healthy
# 2. Registers a test user with a unique email
# 3. Extracts and exports INFERADB_SESSION_TOKEN
# 4. Works for both local use and GitHub Actions
#
# Usage: ./scripts/bootstrap-test-env.sh [OPTIONS]
#
# Options:
#   --endpoint URL    Control API endpoint (default: http://localhost:9090)
#   --timeout SECS    Maximum wait time for Control (default: 120)
#   --export-file     Write exports to file for sourcing (default: .env.test)
#   --github-output   Write to GITHUB_OUTPUT for Actions
#   --quiet           Suppress progress output
#   --help            Show this help message
#
# Environment Variables:
#   INFERADB_ENDPOINT         Override default endpoint
#   INFERADB_TEST_USER_NAME   Test user name (default: "Terraform Test User")
#   INFERADB_TEST_PASSWORD    Test user password (default: auto-generated)

set -euo pipefail

# Get script directory for relative paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Default configuration
# BASE_URL is the raw Control server URL (for health checks)
# ENDPOINT includes /control prefix for API calls (since we're hitting Control directly)
# In production, the gateway handles this routing transparently
BASE_URL="${INFERADB_BASE_URL:-http://localhost:9090}"
ENDPOINT="${INFERADB_ENDPOINT:-${BASE_URL}/control}"
TIMEOUT=120
EXPORT_FILE=".env.test"
GITHUB_OUTPUT_MODE=false
QUIET=false

# Test user configuration
TEST_USER_NAME="${INFERADB_TEST_USER_NAME:-Terraform Test User}"
TEST_PASSWORD="${INFERADB_TEST_PASSWORD:-TestPassword123!$(date +%s)}"

# Parse arguments
while [[ $# -gt 0 ]]; do
    case "$1" in
        --endpoint)
            ENDPOINT="$2"
            shift 2
            ;;
        --timeout)
            TIMEOUT="$2"
            shift 2
            ;;
        --export-file)
            EXPORT_FILE="$2"
            shift 2
            ;;
        --github-output)
            GITHUB_OUTPUT_MODE=true
            shift
            ;;
        --quiet)
            QUIET=true
            shift
            ;;
        --help)
            head -30 "$0" | grep '^#' | sed 's/^# \?//'
            exit 0
            ;;
        *)
            echo "Unknown option: $1" >&2
            exit 1
            ;;
    esac
done

# Logging functions - all log to stderr to avoid mixing with function return values
log() {
    if [[ "$QUIET" != "true" ]]; then
        echo "[bootstrap] $*" >&2
    fi
}

log_error() {
    echo "[bootstrap] ERROR: $*" >&2
}

log_success() {
    if [[ "$QUIET" != "true" ]]; then
        echo "[bootstrap] SUCCESS: $*" >&2
    fi
}

# Generate unique test email using timestamp and random suffix
generate_test_email() {
    local timestamp
    local random_suffix
    timestamp=$(date +%s)
    random_suffix=$(head -c 4 /dev/urandom | xxd -p)
    echo "tf-test-${timestamp}-${random_suffix}@test.inferadb.local"
}

# Wait for Control service
wait_for_control() {
    log "Waiting for Control service at ${BASE_URL}..."

    local wait_args=("--endpoint" "$BASE_URL" "--timeout" "$TIMEOUT")
    if [[ "$QUIET" == "true" ]]; then
        wait_args+=("--quiet")
    fi

    if ! "$SCRIPT_DIR/wait-for-control.sh" "${wait_args[@]}"; then
        log_error "Control service did not become healthy within ${TIMEOUT}s"
        return 1
    fi

    log "Control service is ready"
}

# Register test user
register_test_user() {
    local email="$1"
    local password="$2"
    local name="$3"

    log "Registering test user: ${email}"

    # Make registration request and capture headers
    # Use -c to write cookies to a temp file, -D to capture headers
    local cookie_file
    local header_file
    cookie_file=$(mktemp)
    header_file=$(mktemp)

    trap "rm -f '$cookie_file' '$header_file'" RETURN

    local http_code
    http_code=$(curl -s -o /dev/null -w "%{http_code}" \
        -X POST \
        -H "Content-Type: application/json" \
        -c "$cookie_file" \
        -D "$header_file" \
        -d "{\"email\": \"${email}\", \"password\": \"${password}\", \"name\": \"${name}\"}" \
        "${BASE_URL}/control/v1/auth/register" 2>&1) || {
        log_error "Failed to connect to Control service"
        return 1
    }

    # Check for success (200 OK for registration)
    if [[ "$http_code" != "200" ]]; then
        log_error "Registration failed with HTTP ${http_code}"
        log_error "Check Control service logs for details"
        return 1
    fi

    # Extract session token from Set-Cookie header
    # Format: Set-Cookie: infera_session=<token>; ...
    local session_token
    session_token=$(grep -i "infera_session" "$cookie_file" | awk '{print $NF}' | head -1)

    # Fallback: try parsing from Set-Cookie header directly
    if [[ -z "$session_token" ]]; then
        session_token=$(grep -i "Set-Cookie.*infera_session" "$header_file" | sed 's/.*infera_session=\([^;]*\).*/\1/' | head -1)
    fi

    if [[ -z "$session_token" ]]; then
        log_error "No infera_session cookie in registration response"
        log_error "Headers:"
        cat "$header_file" >&2
        return 1
    fi

    log "Successfully registered test user"
    echo "$session_token"
}

# Verify email via MailHog
# This extracts the verification token from MailHog and calls the verify endpoint
verify_email_via_mailhog() {
    local email="$1"
    local mailhog_api="${MAILHOG_API:-http://localhost:8025}"

    log "Verifying email via MailHog..."

    # Wait a moment for email to arrive
    sleep 1

    # Fetch the latest email from MailHog
    local email_body
    email_body=$(curl -sf "${mailhog_api}/api/v2/messages" 2>/dev/null) || {
        log_error "Failed to connect to MailHog at ${mailhog_api}"
        log_error "Is MailHog running? Check: curl ${mailhog_api}"
        return 1
    }

    # Extract verification token using Python (more reliable than sed/awk for JSON + email parsing)
    local token
    token=$(echo "$email_body" | python3 -c "
import json, sys, re, quopri
try:
    msgs = json.load(sys.stdin).get('items', [])
    for msg in msgs:
        body = msg.get('Content', {}).get('Body', '') or ''
        # Decode quoted-printable encoding
        try:
            decoded = quopri.decodestring(body.encode()).decode()
        except:
            decoded = body
        # Look for verification token in URL
        match = re.search(r'verify-email\?token=([a-zA-Z0-9_-]+)', decoded)
        if match:
            print(match.group(1))
            sys.exit(0)
    sys.exit(1)
except Exception as e:
    sys.exit(1)
" 2>/dev/null) || {
        log_error "Failed to extract verification token from email"
        return 1
    }

    if [[ -z "$token" ]]; then
        log_error "No verification token found in emails"
        return 1
    fi

    # Call the verification endpoint
    local http_code
    http_code=$(curl -s -o /dev/null -w "%{http_code}" \
        -X POST \
        -H "Content-Type: application/json" \
        -d "{\"token\":\"${token}\"}" \
        "${BASE_URL}/control/v1/auth/verify-email" 2>&1) || {
        log_error "Failed to call verification endpoint"
        return 1
    }

    if [[ "$http_code" != "200" ]]; then
        log_error "Email verification failed with HTTP ${http_code}"
        return 1
    fi

    log "Email verified successfully"
}

# Export credentials
export_credentials() {
    local endpoint="$1"
    local session_token="$2"
    local email="$3"

    # Export to environment (for current shell)
    export INFERADB_ENDPOINT="$endpoint"
    export INFERADB_SESSION_TOKEN="$session_token"

    # Write to export file for sourcing
    if [[ -n "$EXPORT_FILE" ]]; then
        cat > "$EXPORT_FILE" <<EOF
# InferaDB Test Environment Credentials
# Generated by bootstrap-test-env.sh at $(date -u +"%Y-%m-%dT%H:%M:%SZ")
# DO NOT COMMIT THIS FILE

export INFERADB_ENDPOINT="${endpoint}"
export INFERADB_SESSION_TOKEN="${session_token}"
export INFERADB_TEST_EMAIL="${email}"
EOF
        log "Credentials written to ${EXPORT_FILE}"
        log "Source with: source ${EXPORT_FILE}"
    fi

    # Write to GITHUB_OUTPUT for GitHub Actions
    if [[ "$GITHUB_OUTPUT_MODE" == "true" ]]; then
        if [[ -n "${GITHUB_OUTPUT:-}" ]]; then
            {
                echo "inferadb_endpoint=${endpoint}"
                echo "inferadb_session_token=${session_token}"
                echo "inferadb_test_email=${email}"
            } >> "$GITHUB_OUTPUT"
            log "Credentials written to GITHUB_OUTPUT"
        else
            log_error "GITHUB_OUTPUT environment variable not set"
            return 1
        fi
    fi

    # Also set GITHUB_ENV for subsequent steps
    if [[ "$GITHUB_OUTPUT_MODE" == "true" ]] && [[ -n "${GITHUB_ENV:-}" ]]; then
        {
            echo "INFERADB_ENDPOINT=${endpoint}"
            echo "INFERADB_SESSION_TOKEN=${session_token}"
        } >> "$GITHUB_ENV"
        log "Environment variables written to GITHUB_ENV"
    fi
}

# Main execution
main() {
    log "Starting test environment bootstrap"
    log "Endpoint: ${ENDPOINT}"

    # Step 1: Wait for Control service
    wait_for_control || exit 1

    # Step 2: Generate unique test credentials
    local test_email
    test_email=$(generate_test_email)
    log "Generated test email: ${test_email}"

    # Step 3: Register test user
    local session_token
    session_token=$(register_test_user "$test_email" "$TEST_PASSWORD" "$TEST_USER_NAME") || exit 1

    # Step 4: Verify email (required before creating organizations)
    verify_email_via_mailhog "$test_email" || exit 1

    # Step 5: Export credentials
    export_credentials "$ENDPOINT" "$session_token" "$test_email" || exit 1

    log_success "Test environment bootstrap complete"
    log_success "INFERADB_ENDPOINT=${ENDPOINT}"
    log_success "INFERADB_SESSION_TOKEN=${session_token:0:20}..."

    # Print usage hint for local development
    if [[ "$GITHUB_OUTPUT_MODE" != "true" ]] && [[ -n "$EXPORT_FILE" ]]; then
        echo ""
        echo "To use these credentials:"
        echo "  source ${EXPORT_FILE}"
        echo "  make testacc"
    fi
}

main "$@"
