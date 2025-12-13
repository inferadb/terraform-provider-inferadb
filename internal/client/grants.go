// Copyright 2025 InferaDB
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"fmt"
)

// CreateVaultUserGrant grants a user access to a vault.
func (c *Client) CreateVaultUserGrant(ctx context.Context, orgID, vaultID string, req CreateVaultUserGrantRequest) (*VaultUserGrant, error) {
	var grant VaultUserGrant
	if err := c.post(ctx, fmt.Sprintf("/v1/organizations/%s/vaults/%s/user-grants", orgID, vaultID), req, &grant); err != nil {
		return nil, fmt.Errorf("failed to create user grant: %w", err)
	}
	return &grant, nil
}

// GetVaultUserGrant retrieves a user grant by ID.
func (c *Client) GetVaultUserGrant(ctx context.Context, orgID, vaultID, grantID string) (*VaultUserGrant, error) {
	var grant VaultUserGrant
	if err := c.get(ctx, fmt.Sprintf("/v1/organizations/%s/vaults/%s/user-grants/%s", orgID, vaultID, grantID), &grant); err != nil {
		return nil, fmt.Errorf("failed to get user grant: %w", err)
	}
	return &grant, nil
}

// UpdateVaultUserGrant updates a user grant's role.
func (c *Client) UpdateVaultUserGrant(ctx context.Context, orgID, vaultID, grantID string, req UpdateVaultUserGrantRequest) (*VaultUserGrant, error) {
	var grant VaultUserGrant
	if err := c.patch(ctx, fmt.Sprintf("/v1/organizations/%s/vaults/%s/user-grants/%s", orgID, vaultID, grantID), req, &grant); err != nil {
		return nil, fmt.Errorf("failed to update user grant: %w", err)
	}
	return &grant, nil
}

// DeleteVaultUserGrant revokes a user's access to a vault.
func (c *Client) DeleteVaultUserGrant(ctx context.Context, orgID, vaultID, grantID string) error {
	if err := c.delete(ctx, fmt.Sprintf("/v1/organizations/%s/vaults/%s/user-grants/%s", orgID, vaultID, grantID)); err != nil {
		return fmt.Errorf("failed to delete user grant: %w", err)
	}
	return nil
}

// ListVaultUserGrants lists all user grants for a vault.
func (c *Client) ListVaultUserGrants(ctx context.Context, orgID, vaultID string) ([]VaultUserGrant, error) {
	var grants []VaultUserGrant
	if err := c.get(ctx, fmt.Sprintf("/v1/organizations/%s/vaults/%s/user-grants", orgID, vaultID), &grants); err != nil {
		return nil, fmt.Errorf("failed to list user grants: %w", err)
	}
	return grants, nil
}

// CreateVaultTeamGrant grants a team access to a vault.
func (c *Client) CreateVaultTeamGrant(ctx context.Context, orgID, vaultID string, req CreateVaultTeamGrantRequest) (*VaultTeamGrant, error) {
	var grant VaultTeamGrant
	if err := c.post(ctx, fmt.Sprintf("/v1/organizations/%s/vaults/%s/team-grants", orgID, vaultID), req, &grant); err != nil {
		return nil, fmt.Errorf("failed to create team grant: %w", err)
	}
	return &grant, nil
}

// GetVaultTeamGrant retrieves a team grant by ID.
func (c *Client) GetVaultTeamGrant(ctx context.Context, orgID, vaultID, grantID string) (*VaultTeamGrant, error) {
	var grant VaultTeamGrant
	if err := c.get(ctx, fmt.Sprintf("/v1/organizations/%s/vaults/%s/team-grants/%s", orgID, vaultID, grantID), &grant); err != nil {
		return nil, fmt.Errorf("failed to get team grant: %w", err)
	}
	return &grant, nil
}

// UpdateVaultTeamGrant updates a team grant's role.
func (c *Client) UpdateVaultTeamGrant(ctx context.Context, orgID, vaultID, grantID string, req UpdateVaultTeamGrantRequest) (*VaultTeamGrant, error) {
	var grant VaultTeamGrant
	if err := c.patch(ctx, fmt.Sprintf("/v1/organizations/%s/vaults/%s/team-grants/%s", orgID, vaultID, grantID), req, &grant); err != nil {
		return nil, fmt.Errorf("failed to update team grant: %w", err)
	}
	return &grant, nil
}

// DeleteVaultTeamGrant revokes a team's access to a vault.
func (c *Client) DeleteVaultTeamGrant(ctx context.Context, orgID, vaultID, grantID string) error {
	if err := c.delete(ctx, fmt.Sprintf("/v1/organizations/%s/vaults/%s/team-grants/%s", orgID, vaultID, grantID)); err != nil {
		return fmt.Errorf("failed to delete team grant: %w", err)
	}
	return nil
}

// ListVaultTeamGrants lists all team grants for a vault.
func (c *Client) ListVaultTeamGrants(ctx context.Context, orgID, vaultID string) ([]VaultTeamGrant, error) {
	var grants []VaultTeamGrant
	if err := c.get(ctx, fmt.Sprintf("/v1/organizations/%s/vaults/%s/team-grants", orgID, vaultID), &grants); err != nil {
		return nil, fmt.Errorf("failed to list team grants: %w", err)
	}
	return grants, nil
}
