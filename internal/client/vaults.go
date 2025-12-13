// Copyright 2025 InferaDB
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"fmt"
)

// CreateVault creates a new vault in an organization.
func (c *Client) CreateVault(ctx context.Context, orgID string, req CreateVaultRequest) (*Vault, error) {
	var resp VaultResponse
	if err := c.post(ctx, fmt.Sprintf("/v1/organizations/%s/vaults", orgID), req, &resp); err != nil {
		return nil, fmt.Errorf("failed to create vault: %w", err)
	}
	return &resp.Vault, nil
}

// GetVault retrieves a vault by ID.
func (c *Client) GetVault(ctx context.Context, orgID, vaultID string) (*Vault, error) {
	// Note: GET endpoint returns vault directly, not wrapped
	var vault Vault
	if err := c.get(ctx, fmt.Sprintf("/v1/organizations/%s/vaults/%s", orgID, vaultID), &vault); err != nil {
		return nil, fmt.Errorf("failed to get vault: %w", err)
	}
	return &vault, nil
}

// UpdateVault updates a vault.
func (c *Client) UpdateVault(ctx context.Context, orgID, vaultID string, req UpdateVaultRequest) (*Vault, error) {
	var resp VaultResponse
	if err := c.patch(ctx, fmt.Sprintf("/v1/organizations/%s/vaults/%s", orgID, vaultID), req, &resp); err != nil {
		return nil, fmt.Errorf("failed to update vault: %w", err)
	}
	return &resp.Vault, nil
}

// DeleteVault deletes a vault.
func (c *Client) DeleteVault(ctx context.Context, orgID, vaultID string) error {
	if err := c.delete(ctx, fmt.Sprintf("/v1/organizations/%s/vaults/%s", orgID, vaultID)); err != nil {
		return fmt.Errorf("failed to delete vault: %w", err)
	}
	return nil
}

// ListVaults lists all vaults in an organization.
func (c *Client) ListVaults(ctx context.Context, orgID string) ([]Vault, error) {
	var resp VaultListResponse
	if err := c.get(ctx, fmt.Sprintf("/v1/organizations/%s/vaults", orgID), &resp); err != nil {
		return nil, fmt.Errorf("failed to list vaults: %w", err)
	}
	return resp.Vaults, nil
}
