// Copyright 2025 InferaDB
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"fmt"
)

// CreateCertificate generates a new Ed25519 certificate for a client.
// The private key is returned ONLY in the response to this call and cannot be retrieved later.
func (c *Client) CreateCertificate(ctx context.Context, orgID, clientID string, req CreateCertificateRequest) (*ClientCertificate, error) {
	var cert ClientCertificate
	if err := c.post(ctx, fmt.Sprintf("/v1/organizations/%s/clients/%s/certificates", orgID, clientID), req, &cert); err != nil {
		return nil, fmt.Errorf("failed to create certificate: %w", err)
	}
	return &cert, nil
}

// GetCertificate retrieves a certificate by ID (without the private key).
func (c *Client) GetCertificate(ctx context.Context, orgID, clientID, certID string) (*ClientCertificate, error) {
	var cert ClientCertificate
	if err := c.get(ctx, fmt.Sprintf("/v1/organizations/%s/clients/%s/certificates/%s", orgID, clientID, certID), &cert); err != nil {
		return nil, fmt.Errorf("failed to get certificate: %w", err)
	}
	return &cert, nil
}

// DeleteCertificate permanently deletes a certificate.
func (c *Client) DeleteCertificate(ctx context.Context, orgID, clientID, certID string) error {
	if err := c.delete(ctx, fmt.Sprintf("/v1/organizations/%s/clients/%s/certificates/%s", orgID, clientID, certID)); err != nil {
		return fmt.Errorf("failed to delete certificate: %w", err)
	}
	return nil
}

// RevokeCertificate revokes a certificate, immediately invalidating all tokens issued with it.
func (c *Client) RevokeCertificate(ctx context.Context, orgID, clientID, certID string) error {
	if err := c.post(ctx, fmt.Sprintf("/v1/organizations/%s/clients/%s/certificates/%s/revoke", orgID, clientID, certID), nil, nil); err != nil {
		return fmt.Errorf("failed to revoke certificate: %w", err)
	}
	return nil
}

// ListCertificates lists all certificates for a client.
func (c *Client) ListCertificates(ctx context.Context, orgID, clientID string) ([]ClientCertificate, error) {
	var certs []ClientCertificate
	if err := c.get(ctx, fmt.Sprintf("/v1/organizations/%s/clients/%s/certificates", orgID, clientID), &certs); err != nil {
		return nil, fmt.Errorf("failed to list certificates: %w", err)
	}
	return certs, nil
}
