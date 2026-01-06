#!/usr/bin/env bash
# Copyright 2025 InferaDB
# SPDX-License-Identifier: Apache-2.0
#
# wait-for-control.sh - Wait for the Control service to be healthy
#
# Usage: ./scripts/wait-for-control.sh [OPTIONS]
#
# Options:
#   --endpoint URL    Control API endpoint (default: http://localhost:9090)
#   --timeout SECS    Maximum wait time in seconds (default: 120)
#   --interval SECS   Polling interval in seconds (default: 2)
#   --quiet           Suppress progress output
#   --help            Show this help message

set -euo pipefail

# Default configuration
ENDPOINT="${INFERADB_ENDPOINT:-http://localhost:9090}"
TIMEOUT=120
INTERVAL=2
QUIET=false

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
        --interval)
            INTERVAL="$2"
            shift 2
            ;;
        --quiet)
            QUIET=true
            shift
            ;;
        --help)
            head -20 "$0" | grep '^#' | sed 's/^# \?//'
            exit 0
            ;;
        *)
            echo "Unknown option: $1" >&2
            exit 1
            ;;
    esac
done

# Logging functions
log() {
    if [[ "$QUIET" != "true" ]]; then
        echo "[wait-for-control] $*"
    fi
}

log_error() {
    echo "[wait-for-control] ERROR: $*" >&2
}

# Health check URL (Kubernetes-style health endpoint)
HEALTH_URL="${ENDPOINT}/healthz"

log "Waiting for Control service at ${ENDPOINT}"
log "Health check URL: ${HEALTH_URL}"
log "Timeout: ${TIMEOUT}s, Interval: ${INTERVAL}s"

start_time=$(date +%s)
attempts=0

while true; do
    attempts=$((attempts + 1))
    current_time=$(date +%s)
    elapsed=$((current_time - start_time))

    # Check timeout
    if [[ $elapsed -ge $TIMEOUT ]]; then
        log_error "Timeout after ${TIMEOUT}s waiting for Control service"
        log_error "Last attempt was #${attempts}"
        exit 1
    fi

    # Perform health check
    if curl -sf --max-time 5 "${HEALTH_URL}" > /dev/null 2>&1; then
        log "Control service is healthy (attempt #${attempts}, ${elapsed}s elapsed)"
        exit 0
    fi

    # Log progress
    remaining=$((TIMEOUT - elapsed))
    log "Attempt #${attempts} failed, retrying in ${INTERVAL}s (${remaining}s remaining)"

    sleep "$INTERVAL"
done
