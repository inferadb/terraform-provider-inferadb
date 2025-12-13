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
var _ resource.Resource = &TeamMemberResource{}
var _ resource.ResourceWithImportState = &TeamMemberResource{}

// NewTeamMemberResource creates a new team member resource.
func NewTeamMemberResource() resource.Resource {
	return &TeamMemberResource{}
}

// TeamMemberResource defines the resource implementation.
type TeamMemberResource struct {
	client *client.Client
}

// TeamMemberResourceModel describes the resource data model.
type TeamMemberResourceModel struct {
	ID             types.String `tfsdk:"id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	TeamID         types.String `tfsdk:"team_id"`
	UserID         types.String `tfsdk:"user_id"`
	Role           types.String `tfsdk:"role"`
	CreatedAt      types.String `tfsdk:"created_at"`
}

// Metadata sets the resource type name.
func (r *TeamMemberResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_member"
}

// Schema defines the resource schema.
func (r *TeamMemberResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Manages a team member in InferaDB. Team members represent user memberships in teams with specific roles.

## Example Usage

` + "```hcl" + `
resource "inferadb_team_member" "example" {
  organization_id = inferadb_organization.example.id
  team_id         = inferadb_team.dev.id
  user_id         = "user_123456789"
  role            = "member"
}

resource "inferadb_team_member" "maintainer" {
  organization_id = inferadb_organization.example.id
  team_id         = inferadb_team.dev.id
  user_id         = "user_987654321"
  role            = "maintainer"
}
` + "```" + `

## Import

Team members can be imported using the format: ` + "`<org_id>/<team_id>/<member_id>`" + `

` + "```shell" + `
terraform import inferadb_team_member.example org_123/team_456/member_789
` + "```",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Membership ID",
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
			"team_id": schema.StringAttribute{
				MarkdownDescription: "Parent team ID",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user_id": schema.StringAttribute{
				MarkdownDescription: "User ID to add",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "Role: maintainer, member",
				Required:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "ISO 8601 timestamp when membership was created",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *TeamMemberResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a new team member.
func (r *TeamMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TeamMemberResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the team member
	member, err := r.client.AddTeamMember(
		ctx,
		data.OrganizationID.ValueString(),
		data.TeamID.ValueString(),
		client.AddTeamMemberRequest{
			UserID: data.UserID.ValueString(),
			Role:   data.Role.ValueString(),
		},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating team member",
			fmt.Sprintf("Could not add team member: %s", err.Error()),
		)
		return
	}

	// Map response to model
	data.ID = types.StringValue(member.ID)
	data.CreatedAt = types.StringValue(member.CreatedAt)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the current team member state.
func (r *TeamMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TeamMemberResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the team member
	member, err := r.client.GetTeamMember(
		ctx,
		data.OrganizationID.ValueString(),
		data.TeamID.ValueString(),
		data.ID.ValueString(),
	)
	if err != nil {
		// Handle 404 - resource was deleted outside Terraform
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading team member",
			fmt.Sprintf("Could not read team member %s: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	// Map response to model
	data.UserID = types.StringValue(member.UserID)
	data.Role = types.StringValue(member.Role)
	data.CreatedAt = types.StringValue(member.CreatedAt)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates an existing team member.
func (r *TeamMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TeamMemberResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the team member
	member, err := r.client.UpdateTeamMember(
		ctx,
		data.OrganizationID.ValueString(),
		data.TeamID.ValueString(),
		data.ID.ValueString(),
		client.UpdateTeamMemberRequest{
			Role: data.Role.ValueString(),
		},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating team member",
			fmt.Sprintf("Could not update team member: %s", err.Error()),
		)
		return
	}

	// Map response to model
	data.Role = types.StringValue(member.Role)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes a team member.
func (r *TeamMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TeamMemberResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the team member
	err := r.client.RemoveTeamMember(
		ctx,
		data.OrganizationID.ValueString(),
		data.TeamID.ValueString(),
		data.ID.ValueString(),
	)
	if err != nil {
		// If resource is already gone, that's okay
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			return
		}

		resp.Diagnostics.AddError(
			"Error deleting team member",
			fmt.Sprintf("Could not remove team member: %s", err.Error()),
		)
		return
	}
}

// ImportState imports a team member using the format: <org_id>/<team_id>/<member_id>
func (r *TeamMemberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 3 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in format: <org_id>/<team_id>/<member_id>, got: %s", req.ID),
		)
		return
	}

	orgID := parts[0]
	teamID := parts[1]
	memberID := parts[2]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), orgID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("team_id"), teamID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), memberID)...)
}
