// Copyright 2025 InferaDB
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"fmt"
)

// CreateClient creates a new client (backend service identity) in an organization.
func (c *Client) CreateClient(ctx context.Context, orgID string, req CreateClientRequest) (*InferaClient, error) {
	var resp ClientResponse
	if err := c.post(ctx, fmt.Sprintf("/v1/organizations/%s/clients", orgID), req, &resp); err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	return &resp.Client, nil
}

// GetClient retrieves a client by ID.
func (c *Client) GetClient(ctx context.Context, orgID, clientID string) (*InferaClient, error) {
	var resp ClientResponse
	if err := c.get(ctx, fmt.Sprintf("/v1/organizations/%s/clients/%s", orgID, clientID), &resp); err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}
	return &resp.Client, nil
}

// UpdateClient updates a client.
func (c *Client) UpdateClient(ctx context.Context, orgID, clientID string, req UpdateClientRequest) (*InferaClient, error) {
	var resp ClientResponse
	if err := c.patch(ctx, fmt.Sprintf("/v1/organizations/%s/clients/%s", orgID, clientID), req, &resp); err != nil {
		return nil, fmt.Errorf("failed to update client: %w", err)
	}
	return &resp.Client, nil
}

// DeleteClient deletes a client.
func (c *Client) DeleteClient(ctx context.Context, orgID, clientID string) error {
	if err := c.delete(ctx, fmt.Sprintf("/v1/organizations/%s/clients/%s", orgID, clientID)); err != nil {
		return fmt.Errorf("failed to delete client: %w", err)
	}
	return nil
}

// DeactivateClient deactivates a client, revoking all certificates and tokens.
func (c *Client) DeactivateClient(ctx context.Context, orgID, clientID string) error {
	if err := c.post(ctx, fmt.Sprintf("/v1/organizations/%s/clients/%s/deactivate", orgID, clientID), nil, nil); err != nil {
		return fmt.Errorf("failed to deactivate client: %w", err)
	}
	return nil
}

// ListClients lists all clients in an organization.
func (c *Client) ListClients(ctx context.Context, orgID string) ([]InferaClient, error) {
	var resp ClientListResponse
	if err := c.get(ctx, fmt.Sprintf("/v1/organizations/%s/clients", orgID), &resp); err != nil {
		return nil, fmt.Errorf("failed to list clients: %w", err)
	}
	return resp.Clients, nil
}
