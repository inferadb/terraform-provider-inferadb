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
var _ resource.Resource = &VaultUserGrantResource{}
var _ resource.ResourceWithImportState = &VaultUserGrantResource{}

// NewVaultUserGrantResource creates a new vault user grant resource.
func NewVaultUserGrantResource() resource.Resource {
	return &VaultUserGrantResource{}
}

// VaultUserGrantResource defines the resource implementation.
type VaultUserGrantResource struct {
	client *client.Client
}

// VaultUserGrantResourceModel describes the resource data model.
type VaultUserGrantResourceModel struct {
	ID              types.String `tfsdk:"id"`
	OrganizationID  types.String `tfsdk:"organization_id"`
	VaultID         types.String `tfsdk:"vault_id"`
	UserID          types.String `tfsdk:"user_id"`
	Role            types.String `tfsdk:"role"`
	GrantedAt       types.String `tfsdk:"granted_at"`
	GrantedByUserID types.String `tfsdk:"granted_by_user_id"`
}

// Metadata sets the resource type name.
func (r *VaultUserGrantResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vault_user_grant"
}

// Schema defines the resource schema.
func (r *VaultUserGrantResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a vault user grant in InferaDB. Grants specify user permissions for a vault.",

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
			"user_id": schema.StringAttribute{
				MarkdownDescription: "User ID",
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
func (r *VaultUserGrantResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = c
}

// Create creates a new vault user grant.
func (r *VaultUserGrantResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VaultUserGrantResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	grant, err := r.client.CreateVaultUserGrant(ctx,
		data.OrganizationID.ValueString(),
		data.VaultID.ValueString(),
		client.CreateVaultUserGrantRequest{
			UserID: data.UserID.ValueString(),
			Role:   data.Role.ValueString(),
		},
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create vault user grant: %s", err))
		return
	}

	data.ID = types.StringValue(grant.ID)
	data.GrantedAt = types.StringValue(grant.GrantedAt)
	data.GrantedByUserID = types.StringValue(grant.GrantedByUserID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the current vault user grant state.
func (r *VaultUserGrantResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VaultUserGrantResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	grant, err := r.client.GetVaultUserGrant(ctx,
		data.OrganizationID.ValueString(),
		data.VaultID.ValueString(),
		data.ID.ValueString(),
	)
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read vault user grant: %s", err))
		return
	}

	data.UserID = types.StringValue(grant.UserID)
	data.Role = types.StringValue(grant.Role)
	data.GrantedAt = types.StringValue(grant.GrantedAt)
	data.GrantedByUserID = types.StringValue(grant.GrantedByUserID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates an existing vault user grant.
func (r *VaultUserGrantResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data VaultUserGrantResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	grant, err := r.client.UpdateVaultUserGrant(ctx,
		data.OrganizationID.ValueString(),
		data.VaultID.ValueString(),
		data.ID.ValueString(),
		client.UpdateVaultUserGrantRequest{
			Role: data.Role.ValueString(),
		},
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update vault user grant: %s", err))
		return
	}

	data.Role = types.StringValue(grant.Role)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes a vault user grant.
func (r *VaultUserGrantResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VaultUserGrantResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteVaultUserGrant(ctx,
		data.OrganizationID.ValueString(),
		data.VaultID.ValueString(),
		data.ID.ValueString(),
	)
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete vault user grant: %s", err))
		return
	}
}

// ImportState imports a vault user grant using the format: <org_id>/<vault_id>/<grant_id>
func (r *VaultUserGrantResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 3 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in format: <org_id>/<vault_id>/<grant_id>, got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vault_id"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[2])...)
}
