// Copyright 2025 InferaDB
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// SnowflakeID is a custom type that can unmarshal both string and number JSON values.
// The InferaDB API returns IDs as numbers, but we want to use them as strings in Terraform.
type SnowflakeID string

// UnmarshalJSON implements json.Unmarshaler for SnowflakeID.
func (s *SnowflakeID) UnmarshalJSON(data []byte) error {
	// Try unmarshaling as string first
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*s = SnowflakeID(str)
		return nil
	}

	// Try unmarshaling as number
	var num int64
	if err := json.Unmarshal(data, &num); err == nil {
		*s = SnowflakeID(strconv.FormatInt(num, 10))
		return nil
	}

	return fmt.Errorf("cannot unmarshal %s into SnowflakeID", string(data))
}

// String returns the string representation of the SnowflakeID.
func (s SnowflakeID) String() string {
	return string(s)
}

// Organization represents an InferaDB organization.
type Organization struct {
	ID          SnowflakeID `json:"id"`
	Name        string      `json:"name"`
	Tier        string      `json:"tier"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at,omitempty"`
	DeletedAt   *string     `json:"deleted_at,omitempty"`
	SuspendedAt *string     `json:"suspended_at,omitempty"`
	Role        string      `json:"role,omitempty"` // User's role in this organization
}

// OrganizationResponse wraps an organization in API responses.
type OrganizationResponse struct {
	Organization Organization `json:"organization"`
}

// OrganizationListResponse wraps a list of organizations in API responses.
type OrganizationListResponse struct {
	Organizations []Organization `json:"organizations"`
	Pagination    Pagination     `json:"pagination"`
}

// Pagination represents API pagination metadata.
type Pagination struct {
	Total   int  `json:"total"`
	Count   int  `json:"count"`
	Offset  int  `json:"offset"`
	Limit   int  `json:"limit"`
	HasMore bool `json:"has_more"`
}

// CreateOrganizationRequest is the request body for creating an organization.
type CreateOrganizationRequest struct {
	Name string `json:"name"`
	Tier string `json:"tier,omitempty"`
}

// UpdateOrganizationRequest is the request body for updating an organization.
type UpdateOrganizationRequest struct {
	Name string `json:"name,omitempty"`
	Tier string `json:"tier,omitempty"`
}

// Vault represents an InferaDB vault.
type Vault struct {
	ID             SnowflakeID `json:"id"`
	OrganizationID SnowflakeID `json:"organization_id"`
	Name           string      `json:"name"`
	Description    string      `json:"description,omitempty"`
	SyncStatus     string      `json:"sync_status,omitempty"`
	SyncError      *string     `json:"sync_error,omitempty"`
	CreatedAt      string      `json:"created_at"`
	UpdatedAt      string      `json:"updated_at,omitempty"`
	DeletedAt      *string     `json:"deleted_at,omitempty"`
}

// VaultResponse wraps a vault in API responses.
type VaultResponse struct {
	Vault Vault `json:"vault"`
}

// VaultListResponse wraps a list of vaults in API responses.
type VaultListResponse struct {
	Vaults     []Vault    `json:"vaults"`
	Pagination Pagination `json:"pagination"`
}

// CreateVaultRequest is the request body for creating a vault.
type CreateVaultRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// UpdateVaultRequest is the request body for updating a vault.
type UpdateVaultRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// InferaClient represents an InferaDB client (backend service identity).
type InferaClient struct {
	ID             SnowflakeID `json:"id"`
	OrganizationID SnowflakeID `json:"organization_id"`
	VaultID        SnowflakeID `json:"vault_id"`
	Name           string      `json:"name"`
	Description    string      `json:"description,omitempty"`
	IsActive       bool        `json:"is_active"`
	CreatedAt      string      `json:"created_at"`
	UpdatedAt      string      `json:"updated_at,omitempty"`
	DeletedAt      *string     `json:"deleted_at,omitempty"`
}

// ClientResponse wraps a client in API responses.
type ClientResponse struct {
	Client InferaClient `json:"client"`
}

// ClientListResponse wraps a list of clients in API responses.
type ClientListResponse struct {
	Clients    []InferaClient `json:"clients"`
	Pagination Pagination     `json:"pagination"`
}

// CreateClientRequest is the request body for creating a client.
type CreateClientRequest struct {
	Name    string `json:"name"`
	VaultID string `json:"vault_id"`
}

// UpdateClientRequest is the request body for updating a client.
type UpdateClientRequest struct {
	Name    string `json:"name,omitempty"`
	VaultID string `json:"vault_id,omitempty"`
}

// ClientCertificate represents an InferaDB client certificate.
type ClientCertificate struct {
	ID              string  `json:"id"`
	ClientID        string  `json:"client_id"`
	Name            string  `json:"name"`
	KID             string  `json:"kid"`
	PublicKeyPEM    string  `json:"public_key_pem,omitempty"`
	PrivateKeyPEM   string  `json:"private_key_pem,omitempty"` // Only returned on creation
	IsActive        bool    `json:"is_active"`
	CreatedAt       string  `json:"created_at"`
	RevokedAt       *string `json:"revoked_at,omitempty"`
	RevokedByUserID *string `json:"revoked_by_user_id,omitempty"`
	DeletedAt       *string `json:"deleted_at,omitempty"`
}

// CreateCertificateRequest is the request body for creating a certificate.
type CreateCertificateRequest struct {
	Name string `json:"name"`
}

// Team represents an InferaDB team.
type Team struct {
	ID             SnowflakeID `json:"id"`
	OrganizationID SnowflakeID `json:"organization_id"`
	Name           string      `json:"name"`
	Description    string      `json:"description,omitempty"`
	CreatedAt      string      `json:"created_at"`
	UpdatedAt      string      `json:"updated_at,omitempty"`
	DeletedAt      *string     `json:"deleted_at,omitempty"`
}

// TeamResponse wraps a team in API responses.
type TeamResponse struct {
	Team Team `json:"team"`
}

// TeamListResponse wraps a list of teams in API responses.
type TeamListResponse struct {
	Teams      []Team     `json:"teams"`
	Pagination Pagination `json:"pagination"`
}

// CreateTeamRequest is the request body for creating a team.
type CreateTeamRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// UpdateTeamRequest is the request body for updating a team.
type UpdateTeamRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// TeamMember represents a team membership.
type TeamMember struct {
	ID        string `json:"id"`
	TeamID    string `json:"team_id"`
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

// AddTeamMemberRequest is the request body for adding a team member.
type AddTeamMemberRequest struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

// UpdateTeamMemberRequest is the request body for updating a team member.
type UpdateTeamMemberRequest struct {
	Role string `json:"role"`
}

// VaultUserGrant represents a user's access grant to a vault.
type VaultUserGrant struct {
	ID              string `json:"id"`
	VaultID         string `json:"vault_id"`
	UserID          string `json:"user_id"`
	Role            string `json:"role"`
	GrantedAt       string `json:"granted_at"`
	GrantedByUserID string `json:"granted_by_user_id"`
}

// CreateVaultUserGrantRequest is the request body for creating a user grant.
type CreateVaultUserGrantRequest struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

// UpdateVaultUserGrantRequest is the request body for updating a user grant.
type UpdateVaultUserGrantRequest struct {
	Role string `json:"role"`
}

// VaultTeamGrant represents a team's access grant to a vault.
type VaultTeamGrant struct {
	ID              string `json:"id"`
	VaultID         string `json:"vault_id"`
	TeamID          string `json:"team_id"`
	Role            string `json:"role"`
	GrantedAt       string `json:"granted_at"`
	GrantedByUserID string `json:"granted_by_user_id"`
}

// CreateVaultTeamGrantRequest is the request body for creating a team grant.
type CreateVaultTeamGrantRequest struct {
	TeamID string `json:"team_id"`
	Role   string `json:"role"`
}

// UpdateVaultTeamGrantRequest is the request body for updating a team grant.
type UpdateVaultTeamGrantRequest struct {
	Role string `json:"role"`
}

// PaginatedResponse represents a paginated API response.
type PaginatedResponse[T any] struct {
	Data       []T `json:"data"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalCount int `json:"total_count"`
	TotalPages int `json:"total_pages"`
}
