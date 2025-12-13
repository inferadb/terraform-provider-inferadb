// Copyright 2025 InferaDB
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/inferadb/terraform-provider-inferadb/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &VaultTeamGrantResource{}
var _ resource.ResourceWithImportState = &VaultTeamGrantResource{}

// NewVaultTeamGrantResource creates a new vault team grant resource.
func NewVaultTeamGrantResource() resource.Resource {
	return &VaultTeamGrantResource{}
}

// VaultTeamGrantResource defines the resource implementation.
type VaultTeamGrantResource struct {
	client *client.Client
}

// VaultTeamGrantResourceModel describes the resource data model.
type VaultTeamGrantResourceModel struct {
	ID              types.String `tfsdk:"id"`
	OrganizationID  types.String `tfsdk:"organization_id"`
	VaultID         types.String `tfsdk:"vault_id"`
	TeamID          types.String `tfsdk:"team_id"`
	Role            types.String `tfsdk:"role"`
	GrantedAt       types.String `tfsdk:"granted_at"`
	GrantedByUserID types.String `tfsdk:"granted_by_user_id"`
}

// Metadata sets the resource type name.
func (r *VaultTeamGrantResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vault_team_grant"
}

// Schema defines the resource schema.
func (r *VaultTeamGrantResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a vault team grant in InferaDB. Grants specify team permissions for a vault.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Grant ID",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Parent organization ID",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vault_id": schema.StringAttribute{
				MarkdownDescription: "Vault ID",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"team_id": schema.StringAttribute{
				MarkdownDescription: "Team ID",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "Role: reader, writer, manager, admin",
				Required:            true,
			},
			"granted_at": schema.StringAttribute{
				MarkdownDescription: "ISO 8601 timestamp when grant was created",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"granted_by_user_id": schema.StringAttribute{
				MarkdownDescription: "ID of user who created the grant",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *VaultTeamGrantResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create creates a new vault team grant.
func (r *VaultTeamGrantResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VaultTeamGrantResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the vault team grant
	grant, err := r.client.CreateVaultTeamGrant(
		ctx,
		data.OrganizationID.ValueString(),
		data.VaultID.ValueString(),
		client.CreateVaultTeamGrantRequest{
			TeamID: data.TeamID.ValueString(),
			Role:   data.Role.ValueString(),
		},
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create vault team grant, got error: %s", err))
		return
	}

	// Map response to model
	data.ID = types.StringValue(grant.ID)
	data.GrantedAt = types.StringValue(grant.GrantedAt)
	data.GrantedByUserID = types.StringValue(grant.GrantedByUserID)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the current vault team grant state.
func (r *VaultTeamGrantResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VaultTeamGrantResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the vault team grant
	grant, err := r.client.GetVaultTeamGrant(
		ctx,
		data.OrganizationID.ValueString(),
		data.VaultID.ValueString(),
		data.ID.ValueString(),
	)
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read vault team grant, got error: %s", err))
		return
	}

	// Map response to model
	data.TeamID = types.StringValue(grant.TeamID)
	data.Role = types.StringValue(grant.Role)
	data.GrantedAt = types.StringValue(grant.GrantedAt)
	data.GrantedByUserID = types.StringValue(grant.GrantedByUserID)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates an existing vault team grant.
func (r *VaultTeamGrantResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data VaultTeamGrantResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the vault team grant
	grant, err := r.client.UpdateVaultTeamGrant(
		ctx,
		data.OrganizationID.ValueString(),
		data.VaultID.ValueString(),
		data.ID.ValueString(),
		client.UpdateVaultTeamGrantRequest{
			Role: data.Role.ValueString(),
		},
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update vault team grant, got error: %s", err))
		return
	}

	// Map response to model
	data.Role = types.StringValue(grant.Role)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes a vault team grant.
func (r *VaultTeamGrantResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VaultTeamGrantResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the vault team grant
	err := r.client.DeleteVaultTeamGrant(
		ctx,
		data.OrganizationID.ValueString(),
		data.VaultID.ValueString(),
		data.ID.ValueString(),
	)
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			// Resource already deleted, treat as success
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete vault team grant, got error: %s", err))
		return
	}
}

// ImportState imports a vault team grant using the format: <org_id>/<vault_id>/<grant_id>
func (r *VaultTeamGrantResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 3 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in format: <org_id>/<vault_id>/<grant_id>, got: %s", req.ID),
		)
		return
	}

	orgID := parts[0]
	vaultID := parts[1]
	grantID := parts[2]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), orgID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vault_id"), vaultID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), grantID)...)
}
