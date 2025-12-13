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
	_ datasource.DataSource              = &ClientDataSource{}
	_ datasource.DataSourceWithConfigure = &ClientDataSource{}
)

// ClientDataSource defines the data source implementation.
type ClientDataSource struct {
	client *client.Client
}

// ClientDataSourceModel describes the data source data model.
type ClientDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	Name           types.String `tfsdk:"name"`
	VaultID        types.String `tfsdk:"vault_id"`
	IsActive       types.Bool   `tfsdk:"is_active"`
	CreatedAt      types.String `tfsdk:"created_at"`
}

// NewClientDataSource is a helper function to simplify the provider implementation.
func NewClientDataSource() datasource.DataSource {
	return &ClientDataSource{}
}

// Metadata returns the data source type name.
func (d *ClientDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_client"
}

// Schema defines the schema for the data source.
func (d *ClientDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Reads an InferaDB client.

Use this data source to retrieve information about an existing client (backend service identity) by its ID.

## Example Usage

` + "```hcl" + `
data "inferadb_client" "example" {
  id              = "5555555555555555555"
  organization_id = "1234567890123456789"
}

output "client_name" {
  value = data.inferadb_client.example.name
}
` + "```",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Snowflake ID of the client.",
				Required:            true,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Snowflake ID of the organization that owns this client.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the client.",
				Computed:            true,
			},
			"vault_id": schema.StringAttribute{
				MarkdownDescription: "Snowflake ID of the vault this client is associated with.",
				Computed:            true,
			},
			"is_active": schema.BoolAttribute{
				MarkdownDescription: "Whether the client is active.",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "ISO 8601 timestamp when the client was created.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *ClientDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *ClientDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ClientDataSourceModel

	// Read configuration
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get client from API
	clientResp, err := d.client.GetClient(ctx, data.OrganizationID.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading client",
			fmt.Sprintf("Could not read client %s: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	// Map response to data model
	data.ID = types.StringValue(clientResp.ID.String())
	data.OrganizationID = types.StringValue(clientResp.OrganizationID.String())
	data.Name = types.StringValue(clientResp.Name)
	data.VaultID = types.StringValue(clientResp.VaultID.String())
	data.IsActive = types.BoolValue(clientResp.IsActive)
	data.CreatedAt = types.StringValue(clientResp.CreatedAt)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
