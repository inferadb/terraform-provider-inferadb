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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/inferadb/terraform-provider-inferadb/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &OrganizationResource{}
	_ resource.ResourceWithImportState = &OrganizationResource{}
)

// OrganizationResource defines the resource implementation.
type OrganizationResource struct {
	client *client.Client
}

// OrganizationResourceModel describes the resource data model.
type OrganizationResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Tier        types.String `tfsdk:"tier"`
	CreatedAt   types.String `tfsdk:"created_at"`
	SuspendedAt types.String `tfsdk:"suspended_at"`
}

// NewOrganizationResource is a helper function to simplify the provider implementation.
func NewOrganizationResource() resource.Resource {
	return &OrganizationResource{}
}

// Metadata returns the resource type name.
func (r *OrganizationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

// Schema defines the schema for the resource.
func (r *OrganizationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Manages an InferaDB organization.

Organizations are the top-level entity in InferaDB. They contain vaults, teams, and users. Each organization has a tier that determines its feature set and resource limits.

## Example Usage

` + "```hcl" + `
resource "inferadb_organization" "example" {
  name = "My Organization"
  tier = "dev"
}

resource "inferadb_organization" "production" {
  name = "Production Org"
  tier = "pro"
}
` + "```" + `

## Import

Organizations can be imported using their ID:

` + "```shell" + `
terraform import inferadb_organization.example <org_id>
` + "```",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Snowflake ID of the organization.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the organization.",
				Required:            true,
			},
			"tier": schema.StringAttribute{
				MarkdownDescription: "Tier of the organization. Valid values are `dev`, `pro`, or `max`. Defaults to `dev`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("dev"),
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "ISO 8601 timestamp when the organization was created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"suspended_at": schema.StringAttribute{
				MarkdownDescription: "ISO 8601 timestamp when the organization was suspended, if applicable.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *OrganizationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *OrganizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OrganizationResourceModel

	// Retrieve values from plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the organization via API
	createReq := client.CreateOrganizationRequest{
		Name: plan.Name.ValueString(),
		Tier: plan.Tier.ValueString(),
	}

	org, err := r.client.CreateOrganization(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating organization",
			fmt.Sprintf("Could not create organization: %s", err.Error()),
		)
		return
	}

	// Map response to resource model
	plan.ID = types.StringValue(org.ID.String())
	plan.Name = types.StringValue(org.Name)
	plan.Tier = types.StringValue(org.Tier)
	plan.CreatedAt = types.StringValue(org.CreatedAt)
	if org.SuspendedAt != nil && *org.SuspendedAt != "" {
		plan.SuspendedAt = types.StringValue(*org.SuspendedAt)
	} else {
		plan.SuspendedAt = types.StringNull()
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *OrganizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OrganizationResourceModel

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get organization from API
	org, err := r.client.GetOrganization(ctx, state.ID.ValueString())
	if err != nil {
		// Handle 404 - resource was deleted outside Terraform
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading organization",
			fmt.Sprintf("Could not read organization %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	// Update state with refreshed data
	state.Name = types.StringValue(org.Name)
	state.Tier = types.StringValue(org.Tier)
	state.CreatedAt = types.StringValue(org.CreatedAt)
	if org.SuspendedAt != nil && *org.SuspendedAt != "" {
		state.SuspendedAt = types.StringValue(*org.SuspendedAt)
	} else {
		state.SuspendedAt = types.StringNull()
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *OrganizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OrganizationResourceModel
	var state OrganizationResourceModel

	// Get current plan and state
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update organization via API
	updateReq := client.UpdateOrganizationRequest{
		Name: plan.Name.ValueString(),
		Tier: plan.Tier.ValueString(),
	}

	org, err := r.client.UpdateOrganization(ctx, state.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating organization",
			fmt.Sprintf("Could not update organization %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	// Update state with response data
	plan.ID = types.StringValue(org.ID.String())
	plan.Name = types.StringValue(org.Name)
	plan.Tier = types.StringValue(org.Tier)
	plan.CreatedAt = types.StringValue(org.CreatedAt)
	if org.SuspendedAt != nil && *org.SuspendedAt != "" {
		plan.SuspendedAt = types.StringValue(*org.SuspendedAt)
	} else {
		plan.SuspendedAt = types.StringNull()
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *OrganizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OrganizationResourceModel

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete organization via API
	err := r.client.DeleteOrganization(ctx, state.ID.ValueString())
	if err != nil {
		// Handle 404 - resource was already deleted
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			return
		}

		resp.Diagnostics.AddError(
			"Error deleting organization",
			fmt.Sprintf("Could not delete organization %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}
}

// ImportState imports the resource into Terraform state.
func (r *OrganizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: <org_id>
	id := strings.TrimSpace(req.ID)
	if id == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Organization ID cannot be empty. Use format: <org_id>",
		)
		return
	}

	// Set the ID and trigger a Read to populate the rest
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
