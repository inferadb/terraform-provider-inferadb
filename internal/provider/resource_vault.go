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
var (
	_ resource.Resource                = &VaultResource{}
	_ resource.ResourceWithImportState = &VaultResource{}
)

// NewVaultResource is a helper function to simplify the provider implementation.
func NewVaultResource() resource.Resource {
	return &VaultResource{}
}

// VaultResource defines the resource implementation.
type VaultResource struct {
	client *client.Client
}

// VaultResourceModel describes the resource data model.
type VaultResourceModel struct {
	ID             types.String `tfsdk:"id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	SyncStatus     types.String `tfsdk:"sync_status"`
	CreatedAt      types.String `tfsdk:"created_at"`
}

// Metadata returns the resource type name.
func (r *VaultResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vault"
}

// Schema defines the schema for the resource.
func (r *VaultResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Manages an InferaDB vault.

Vaults are policy containers within an organization. Each vault contains its own set of
authorization policies, tuples, and client identities, isolated from other vaults.

## Example Usage

` + "```hcl" + `
resource "inferadb_vault" "production" {
  organization_id = inferadb_organization.example.id
  name            = "Production Policies"
  description     = "Authorization policies for production environment"
}
` + "```" + `

## Import

Vaults can be imported using the format ` + "`organization_id/vault_id`" + `:

` + "```" + `
terraform import inferadb_vault.production 123456789/987654321
` + "```",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Snowflake ID of the vault.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "ID of the parent organization.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the vault.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the vault.",
				Optional:            true,
			},
			"sync_status": schema.StringAttribute{
				MarkdownDescription: "Engine sync status. One of: `pending`, `synced`, `failed`.",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "ISO 8601 timestamp of when the vault was created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *VaultResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

// Create creates the resource and sets the initial Terraform state.
func (r *VaultResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan VaultResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the vault via the API
	createReq := client.CreateVaultRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	vault, err := r.client.CreateVault(ctx, plan.OrganizationID.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Vault",
			fmt.Sprintf("Could not create vault: %s", err.Error()),
		)
		return
	}

	// Map response to model
	plan.ID = types.StringValue(vault.ID.String())
	plan.OrganizationID = types.StringValue(vault.OrganizationID.String())
	plan.Name = types.StringValue(vault.Name)
	if vault.Description != "" {
		plan.Description = types.StringValue(vault.Description)
	} else {
		plan.Description = types.StringNull()
	}
	if vault.SyncStatus != "" {
		plan.SyncStatus = types.StringValue(vault.SyncStatus)
	} else {
		plan.SyncStatus = types.StringNull()
	}
	plan.CreatedAt = types.StringValue(vault.CreatedAt)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *VaultResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state VaultResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed vault value from API
	vault, err := r.client.GetVault(ctx, state.OrganizationID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Vault",
			fmt.Sprintf("Could not read vault ID %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	// Map response to model
	state.ID = types.StringValue(vault.ID.String())
	state.OrganizationID = types.StringValue(vault.OrganizationID.String())
	state.Name = types.StringValue(vault.Name)
	if vault.Description != "" {
		state.Description = types.StringValue(vault.Description)
	} else {
		state.Description = types.StringNull()
	}
	if vault.SyncStatus != "" {
		state.SyncStatus = types.StringValue(vault.SyncStatus)
	} else {
		state.SyncStatus = types.StringNull()
	}
	state.CreatedAt = types.StringValue(vault.CreatedAt)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *VaultResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan VaultResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the vault via the API
	updateReq := client.UpdateVaultRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	vault, err := r.client.UpdateVault(ctx, plan.OrganizationID.ValueString(), plan.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Vault",
			fmt.Sprintf("Could not update vault ID %s: %s", plan.ID.ValueString(), err.Error()),
		)
		return
	}

	// Map response to model
	plan.ID = types.StringValue(vault.ID.String())
	plan.OrganizationID = types.StringValue(vault.OrganizationID.String())
	plan.Name = types.StringValue(vault.Name)
	if vault.Description != "" {
		plan.Description = types.StringValue(vault.Description)
	} else {
		plan.Description = types.StringNull()
	}
	if vault.SyncStatus != "" {
		plan.SyncStatus = types.StringValue(vault.SyncStatus)
	} else {
		plan.SyncStatus = types.StringNull()
	}
	plan.CreatedAt = types.StringValue(vault.CreatedAt)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *VaultResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state VaultResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the vault via the API
	err := r.client.DeleteVault(ctx, state.OrganizationID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Vault",
			fmt.Sprintf("Could not delete vault ID %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}
}

// ImportState imports the resource into Terraform state.
func (r *VaultResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expected import ID format: organization_id/vault_id
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID format: organization_id/vault_id, got: %s", req.ID),
		)
		return
	}

	organizationID := parts[0]
	vaultID := parts[1]

	// Set the organization_id and id attributes
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), organizationID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), vaultID)...)
}
