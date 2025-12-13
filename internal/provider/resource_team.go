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

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &TeamResource{}
	_ resource.ResourceWithImportState = &TeamResource{}
)

// TeamResource defines the resource implementation.
type TeamResource struct {
	client *client.Client
}

// TeamResourceModel describes the resource data model.
type TeamResourceModel struct {
	ID             types.String `tfsdk:"id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	CreatedAt      types.String `tfsdk:"created_at"`
}

// NewTeamResource is a helper function to simplify the provider implementation.
func NewTeamResource() resource.Resource {
	return &TeamResource{}
}

// Metadata returns the resource type name.
func (r *TeamResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

// Schema defines the schema for the resource.
func (r *TeamResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Manages an InferaDB team.

Teams are groups of users within an organization. They can be granted access to vaults and used to organize permissions.

## Example Usage

` + "```hcl" + `
resource "inferadb_organization" "example" {
  name = "My Organization"
  tier = "dev"
}

resource "inferadb_team" "engineering" {
  organization_id = inferadb_organization.example.id
  name            = "Engineering"
  description     = "Engineering team with access to production vaults"
}

resource "inferadb_team" "security" {
  organization_id = inferadb_organization.example.id
  name            = "Security Team"
}
` + "```" + `

## Import

Teams can be imported using the format ` + "`<org_id>/<team_id>`" + `:

` + "```shell" + `
terraform import inferadb_team.engineering <org_id>/<team_id>
` + "```",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Snowflake ID of the team.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "ID of the organization this team belongs to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the team.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the team.",
				Optional:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "ISO 8601 timestamp when the team was created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *TeamResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *TeamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TeamResourceModel

	// Retrieve values from plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the team via API
	createReq := client.CreateTeamRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	team, err := r.client.CreateTeam(ctx, plan.OrganizationID.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating team",
			fmt.Sprintf("Could not create team: %s", err.Error()),
		)
		return
	}

	// Map response to resource model
	plan.ID = types.StringValue(team.ID.String())
	plan.OrganizationID = types.StringValue(team.OrganizationID.String())
	plan.Name = types.StringValue(team.Name)
	if team.Description != "" {
		plan.Description = types.StringValue(team.Description)
	} else {
		plan.Description = types.StringNull()
	}
	plan.CreatedAt = types.StringValue(team.CreatedAt)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *TeamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TeamResourceModel

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get team from API
	team, err := r.client.GetTeam(ctx, state.OrganizationID.ValueString(), state.ID.ValueString())
	if err != nil {
		// Handle 404 - resource was deleted outside Terraform
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading team",
			fmt.Sprintf("Could not read team %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	// Update state with refreshed data
	state.OrganizationID = types.StringValue(team.OrganizationID.String())
	state.Name = types.StringValue(team.Name)
	if team.Description != "" {
		state.Description = types.StringValue(team.Description)
	} else {
		state.Description = types.StringNull()
	}
	state.CreatedAt = types.StringValue(team.CreatedAt)

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *TeamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TeamResourceModel
	var state TeamResourceModel

	// Get current plan and state
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update team via API
	updateReq := client.UpdateTeamRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	team, err := r.client.UpdateTeam(ctx, state.OrganizationID.ValueString(), state.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating team",
			fmt.Sprintf("Could not update team %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	// Update state with response data
	plan.ID = types.StringValue(team.ID.String())
	plan.OrganizationID = types.StringValue(team.OrganizationID.String())
	plan.Name = types.StringValue(team.Name)
	if team.Description != "" {
		plan.Description = types.StringValue(team.Description)
	} else {
		plan.Description = types.StringNull()
	}
	plan.CreatedAt = types.StringValue(team.CreatedAt)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *TeamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TeamResourceModel

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete team via API
	err := r.client.DeleteTeam(ctx, state.OrganizationID.ValueString(), state.ID.ValueString())
	if err != nil {
		// Handle 404 - resource was already deleted
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			return
		}

		resp.Diagnostics.AddError(
			"Error deleting team",
			fmt.Sprintf("Could not delete team %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}
}

// ImportState imports the resource into Terraform state.
func (r *TeamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: <org_id>/<team_id>
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Team import ID must be in the format: <org_id>/<team_id>",
		)
		return
	}

	orgID := strings.TrimSpace(parts[0])
	teamID := strings.TrimSpace(parts[1])

	// Set the attributes and trigger a Read to populate the rest
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), orgID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), teamID)...)
}
