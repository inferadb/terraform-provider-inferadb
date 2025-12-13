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
	_ datasource.DataSource              = &VaultDataSource{}
	_ datasource.DataSourceWithConfigure = &VaultDataSource{}
)

// VaultDataSource defines the data source implementation.
type VaultDataSource struct {
	client *client.Client
}

// VaultDataSourceModel describes the data source data model.
type VaultDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	SyncStatus     types.String `tfsdk:"sync_status"`
	CreatedAt      types.String `tfsdk:"created_at"`
}

// NewVaultDataSource is a helper function to simplify the provider implementation.
func NewVaultDataSource() datasource.DataSource {
	return &VaultDataSource{}
}

// Metadata returns the data source type name.
func (d *VaultDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vault"
}

// Schema defines the schema for the data source.
func (d *VaultDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Reads an InferaDB vault.

Use this data source to retrieve information about an existing vault by its ID.

## Example Usage

` + "```hcl" + `
data "inferadb_vault" "example" {
  id              = "9876543210987654321"
  organization_id = "1234567890123456789"
}

output "vault_name" {
  value = data.inferadb_vault.example.name
}
` + "```",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Snowflake ID of the vault.",
				Required:            true,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Snowflake ID of the organization that owns this vault.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the vault.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the vault.",
				Computed:            true,
			},
			"sync_status": schema.StringAttribute{
				MarkdownDescription: "Synchronization status of the vault.",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "ISO 8601 timestamp when the vault was created.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *VaultDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *VaultDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VaultDataSourceModel

	// Read configuration
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get vault from API
	vault, err := d.client.GetVault(ctx, data.OrganizationID.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading vault",
			fmt.Sprintf("Could not read vault %s: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	// Map response to data model
	data.ID = types.StringValue(vault.ID.String())
	data.OrganizationID = types.StringValue(vault.OrganizationID.String())
	data.Name = types.StringValue(vault.Name)
	data.Description = types.StringValue(vault.Description)
	if vault.SyncStatus != "" {
		data.SyncStatus = types.StringValue(vault.SyncStatus)
	} else {
		data.SyncStatus = types.StringNull()
	}
	data.CreatedAt = types.StringValue(vault.CreatedAt)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
