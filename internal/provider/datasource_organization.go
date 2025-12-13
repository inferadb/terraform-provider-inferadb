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
	_ datasource.DataSource              = &OrganizationDataSource{}
	_ datasource.DataSourceWithConfigure = &OrganizationDataSource{}
)

// OrganizationDataSource defines the data source implementation.
type OrganizationDataSource struct {
	client *client.Client
}

// OrganizationDataSourceModel describes the data source data model.
type OrganizationDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Tier        types.String `tfsdk:"tier"`
	CreatedAt   types.String `tfsdk:"created_at"`
	SuspendedAt types.String `tfsdk:"suspended_at"`
}

// NewOrganizationDataSource is a helper function to simplify the provider implementation.
func NewOrganizationDataSource() datasource.DataSource {
	return &OrganizationDataSource{}
}

// Metadata returns the data source type name.
func (d *OrganizationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

// Schema defines the schema for the data source.
func (d *OrganizationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Reads an InferaDB organization.

Use this data source to retrieve information about an existing organization by its ID.

## Example Usage

` + "```hcl" + `
data "inferadb_organization" "example" {
  id = "1234567890123456789"
}

output "org_name" {
  value = data.inferadb_organization.example.name
}
` + "```",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Snowflake ID of the organization.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the organization.",
				Computed:            true,
			},
			"tier": schema.StringAttribute{
				MarkdownDescription: "Tier of the organization (dev, pro, or max).",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "ISO 8601 timestamp when the organization was created.",
				Computed:            true,
			},
			"suspended_at": schema.StringAttribute{
				MarkdownDescription: "ISO 8601 timestamp when the organization was suspended, if applicable.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *OrganizationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *OrganizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationDataSourceModel

	// Read configuration
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get organization from API
	org, err := d.client.GetOrganization(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading organization",
			fmt.Sprintf("Could not read organization %s: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	// Map response to data model
	data.ID = types.StringValue(org.ID.String())
	data.Name = types.StringValue(org.Name)
	data.Tier = types.StringValue(org.Tier)
	data.CreatedAt = types.StringValue(org.CreatedAt)
	if org.SuspendedAt != nil && *org.SuspendedAt != "" {
		data.SuspendedAt = types.StringValue(*org.SuspendedAt)
	} else {
		data.SuspendedAt = types.StringNull()
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
