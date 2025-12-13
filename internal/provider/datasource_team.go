// Copyright 2025 InferaDB
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/inferadb/terraform-provider-inferadb/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &TeamDataSource{}
	_ datasource.DataSourceWithConfigure = &TeamDataSource{}
)

// TeamDataSource defines the data source implementation.
type TeamDataSource struct {
	client *client.Client
}

// TeamDataSourceModel describes the data source data model.
type TeamDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	CreatedAt      types.String `tfsdk:"created_at"`
}

// NewTeamDataSource is a helper function to simplify the provider implementation.
func NewTeamDataSource() datasource.DataSource {
	return &TeamDataSource{}
}

// Metadata returns the data source type name.
func (d *TeamDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

// Schema defines the schema for the data source.
func (d *TeamDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Reads an InferaDB team.

Use this data source to retrieve information about an existing team by its ID.

## Example Usage

` + "```hcl" + `
data "inferadb_team" "example" {
  id              = "7777777777777777777"
  organization_id = "1234567890123456789"
}

output "team_name" {
  value = data.inferadb_team.example.name
}
` + "```",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Snowflake ID of the team.",
				Required:            true,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Snowflake ID of the organization that owns this team.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the team.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the team.",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "ISO 8601 timestamp when the team was created.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *TeamDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

// Read reads the data source's values and updates the state.
func (d *TeamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TeamDataSourceModel

	// Read configuration
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get team from API
	team, err := d.client.GetTeam(ctx, data.OrganizationID.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading team",
			fmt.Sprintf("Could not read team %s: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	// Map response to data model
	data.ID = types.StringValue(team.ID.String())
	data.OrganizationID = types.StringValue(team.OrganizationID.String())
	data.Name = types.StringValue(team.Name)
	data.Description = types.StringValue(team.Description)
	data.CreatedAt = types.StringValue(team.CreatedAt)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
