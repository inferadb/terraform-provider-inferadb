// Copyright 2025 InferaDB
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"fmt"
	"strings"
)

// normalizeTier converts API tier format (TIER_DEV_V1) to Terraform format (dev).
func normalizeTier(apiTier string) string {
	// Map API tiers to user-friendly names
	tierMap := map[string]string{
		"TIER_DEV_V1": "dev",
		"TIER_PRO_V1": "pro",
		"TIER_MAX_V1": "max",
	}
	if tier, ok := tierMap[apiTier]; ok {
		return tier
	}
	// If not found, try to extract from pattern TIER_XXX_V1 -> xxx
	if strings.HasPrefix(apiTier, "TIER_") && strings.HasSuffix(apiTier, "_V1") {
		return strings.ToLower(strings.TrimSuffix(strings.TrimPrefix(apiTier, "TIER_"), "_V1"))
	}
	return apiTier
}

// CreateOrganization creates a new organization.
func (c *Client) CreateOrganization(ctx context.Context, req CreateOrganizationRequest) (*Organization, error) {
	var resp OrganizationResponse
	if err := c.post(ctx, "/v1/organizations", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}
	org := resp.Organization
	// Normalize tier from API format (TIER_DEV_V1) to Terraform format (dev)
	org.Tier = normalizeTier(org.Tier)
	return &org, nil
}

// GetOrganization retrieves an organization by ID.
func (c *Client) GetOrganization(ctx context.Context, id string) (*Organization, error) {
	var resp OrganizationResponse
	if err := c.get(ctx, fmt.Sprintf("/v1/organizations/%s", id), &resp); err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	org := resp.Organization
	org.Tier = normalizeTier(org.Tier)
	return &org, nil
}

// UpdateOrganization updates an organization.
func (c *Client) UpdateOrganization(ctx context.Context, id string, req UpdateOrganizationRequest) (*Organization, error) {
	var resp OrganizationResponse
	if err := c.patch(ctx, fmt.Sprintf("/v1/organizations/%s", id), req, &resp); err != nil {
		return nil, fmt.Errorf("failed to update organization: %w", err)
	}
	org := resp.Organization
	org.Tier = normalizeTier(org.Tier)
	return &org, nil
}

// DeleteOrganization deletes an organization.
func (c *Client) DeleteOrganization(ctx context.Context, id string) error {
	if err := c.delete(ctx, fmt.Sprintf("/v1/organizations/%s", id)); err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}
	return nil
}

// ListOrganizations lists all organizations the user belongs to.
func (c *Client) ListOrganizations(ctx context.Context) ([]Organization, error) {
	var resp OrganizationListResponse
	if err := c.get(ctx, "/v1/organizations", &resp); err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}
	// Normalize tiers
	for i := range resp.Organizations {
		resp.Organizations[i].Tier = normalizeTier(resp.Organizations[i].Tier)
	}
	return resp.Organizations, nil
}
