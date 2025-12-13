// Copyright 2025 InferaDB
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"fmt"
)

// CreateTeam creates a new team in an organization.
func (c *Client) CreateTeam(ctx context.Context, orgID string, req CreateTeamRequest) (*Team, error) {
	var resp TeamResponse
	if err := c.post(ctx, fmt.Sprintf("/v1/organizations/%s/teams", orgID), req, &resp); err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}
	return &resp.Team, nil
}

// GetTeam retrieves a team by ID.
func (c *Client) GetTeam(ctx context.Context, orgID, teamID string) (*Team, error) {
	// Note: GET endpoint returns team directly, not wrapped
	var team Team
	if err := c.get(ctx, fmt.Sprintf("/v1/organizations/%s/teams/%s", orgID, teamID), &team); err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}
	return &team, nil
}

// UpdateTeam updates a team.
func (c *Client) UpdateTeam(ctx context.Context, orgID, teamID string, req UpdateTeamRequest) (*Team, error) {
	var resp TeamResponse
	if err := c.patch(ctx, fmt.Sprintf("/v1/organizations/%s/teams/%s", orgID, teamID), req, &resp); err != nil {
		return nil, fmt.Errorf("failed to update team: %w", err)
	}
	return &resp.Team, nil
}

// DeleteTeam deletes a team.
func (c *Client) DeleteTeam(ctx context.Context, orgID, teamID string) error {
	if err := c.delete(ctx, fmt.Sprintf("/v1/organizations/%s/teams/%s", orgID, teamID)); err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}
	return nil
}

// ListTeams lists all teams in an organization.
func (c *Client) ListTeams(ctx context.Context, orgID string) ([]Team, error) {
	var resp TeamListResponse
	if err := c.get(ctx, fmt.Sprintf("/v1/organizations/%s/teams", orgID), &resp); err != nil {
		return nil, fmt.Errorf("failed to list teams: %w", err)
	}
	return resp.Teams, nil
}

// AddTeamMember adds a user to a team.
func (c *Client) AddTeamMember(ctx context.Context, orgID, teamID string, req AddTeamMemberRequest) (*TeamMember, error) {
	var member TeamMember
	if err := c.post(ctx, fmt.Sprintf("/v1/organizations/%s/teams/%s/members", orgID, teamID), req, &member); err != nil {
		return nil, fmt.Errorf("failed to add team member: %w", err)
	}
	return &member, nil
}

// GetTeamMember retrieves a team member by ID.
func (c *Client) GetTeamMember(ctx context.Context, orgID, teamID, memberID string) (*TeamMember, error) {
	var member TeamMember
	if err := c.get(ctx, fmt.Sprintf("/v1/organizations/%s/teams/%s/members/%s", orgID, teamID, memberID), &member); err != nil {
		return nil, fmt.Errorf("failed to get team member: %w", err)
	}
	return &member, nil
}

// UpdateTeamMember updates a team member's role.
func (c *Client) UpdateTeamMember(ctx context.Context, orgID, teamID, memberID string, req UpdateTeamMemberRequest) (*TeamMember, error) {
	var member TeamMember
	if err := c.patch(ctx, fmt.Sprintf("/v1/organizations/%s/teams/%s/members/%s", orgID, teamID, memberID), req, &member); err != nil {
		return nil, fmt.Errorf("failed to update team member: %w", err)
	}
	return &member, nil
}

// RemoveTeamMember removes a user from a team.
func (c *Client) RemoveTeamMember(ctx context.Context, orgID, teamID, memberID string) error {
	if err := c.delete(ctx, fmt.Sprintf("/v1/organizations/%s/teams/%s/members/%s", orgID, teamID, memberID)); err != nil {
		return fmt.Errorf("failed to remove team member: %w", err)
	}
	return nil
}

// ListTeamMembers lists all members of a team.
func (c *Client) ListTeamMembers(ctx context.Context, orgID, teamID string) ([]TeamMember, error) {
	var members []TeamMember
	if err := c.get(ctx, fmt.Sprintf("/v1/organizations/%s/teams/%s/members", orgID, teamID), &members); err != nil {
		return nil, fmt.Errorf("failed to list team members: %w", err)
	}
	return members, nil
}
