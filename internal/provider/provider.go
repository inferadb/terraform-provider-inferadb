// Copyright 2025 InferaDB
// SPDX-License-Identifier: Apache-2.0

// Package provider implements the InferaDB Terraform provider.
package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/inferadb/terraform-provider-inferadb/internal/client"
)

// Ensure InferaDBProvider satisfies various provider interfaces.
var _ provider.Provider = &InferaDBProvider{}

// InferaDBProvider defines the provider implementation.
type InferaDBProvider struct {
	version string
}

// InferaDBProviderModel describes the provider data model.
type InferaDBProviderModel struct {
	Endpoint     types.String `tfsdk:"endpoint"`
	SessionToken types.String `tfsdk:"session_token"`
}

// New creates a new provider instance.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &InferaDBProvider{
			version: version,
		}
	}
}

// Metadata returns the provider type name.
func (p *InferaDBProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "inferadb"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *InferaDBProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `The InferaDB provider allows you to manage InferaDB Control Plane resources.

## Authentication

The provider supports authentication via session token. You can provide the session token
either in the provider configuration or via the ` + "`INFERADB_SESSION_TOKEN`" + ` environment variable.

To obtain a session token, log in via the InferaDB CLI or web dashboard.

## Example Usage

` + "```hcl" + `
provider "inferadb" {
  endpoint      = "https://api.inferadb.com"
  session_token = var.inferadb_session_token
}

resource "inferadb_organization" "example" {
  name = "My Organization"
  tier = "dev"
}

resource "inferadb_vault" "production" {
  organization_id = inferadb_organization.example.id
  name            = "Production Policies"
}
` + "```" + `
`,
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "InferaDB Control API endpoint. Can also be set via `INFERADB_ENDPOINT` environment variable. Defaults to `https://api.inferadb.com`.",
				Optional:            true,
			},
			"session_token": schema.StringAttribute{
				MarkdownDescription: "Session token for authentication. Can also be set via `INFERADB_SESSION_TOKEN` environment variable. Obtain this by logging in via the InferaDB CLI or web dashboard.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

// Configure prepares an InferaDB API client for data sources and resources.
func (p *InferaDBProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config InferaDBProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Default values with environment variable fallback
	endpoint := "https://api.inferadb.com"
	sessionToken := ""

	if envEndpoint := os.Getenv("INFERADB_ENDPOINT"); envEndpoint != "" {
		endpoint = envEndpoint
	}
	if envToken := os.Getenv("INFERADB_SESSION_TOKEN"); envToken != "" {
		sessionToken = envToken
	}

	// Configuration values override environment variables
	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}
	if !config.SessionToken.IsNull() {
		sessionToken = config.SessionToken.ValueString()
	}

	// Validate required configuration
	if sessionToken == "" {
		resp.Diagnostics.AddError(
			"Missing Session Token",
			"The provider requires a session token for authentication. "+
				"Set the session_token in the provider configuration or "+
				"set the INFERADB_SESSION_TOKEN environment variable.",
		)
		return
	}

	// Create the API client
	apiClient := client.New(client.Config{
		Endpoint:     endpoint,
		SessionToken: sessionToken,
	})

	// Make the client available to data sources and resources
	resp.DataSourceData = apiClient
	resp.ResourceData = apiClient
}

// Resources defines the resources implemented in the provider.
func (p *InferaDBProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewOrganizationResource,
		NewVaultResource,
		NewClientResource,
		NewClientCertificateResource,
		NewTeamResource,
		NewTeamMemberResource,
		NewVaultUserGrantResource,
		NewVaultTeamGrantResource,
	}
}

// DataSources defines the data sources implemented in the provider.
func (p *InferaDBProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewOrganizationDataSource,
		NewVaultDataSource,
		NewClientDataSource,
		NewTeamDataSource,
	}
}
