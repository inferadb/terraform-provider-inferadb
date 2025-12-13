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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/inferadb/terraform-provider-inferadb/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ClientResource{}
var _ resource.ResourceWithImportState = &ClientResource{}

// NewClientResource creates a new client resource.
func NewClientResource() resource.Resource {
	return &ClientResource{}
}

// ClientResource defines the resource implementation.
type ClientResource struct {
	client *client.Client
}

// ClientResourceModel describes the resource data model.
type ClientResourceModel struct {
	ID             types.String `tfsdk:"id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	VaultID        types.String `tfsdk:"vault_id"`
	Name           types.String `tfsdk:"name"`
	IsActive       types.Bool   `tfsdk:"is_active"`
	CreatedAt      types.String `tfsdk:"created_at"`
}

// Metadata sets the resource type name.
func (r *ClientResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_client"
}

// Schema defines the resource schema.
func (r *ClientResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a client (backend service identity) in InferaDB. Clients can authenticate and perform operations on behalf of services.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Snowflake ID",
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
				MarkdownDescription: "Default vault for token generation",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Client name",
				Required:            true,
			},
			"is_active": schema.BoolAttribute{
				MarkdownDescription: "Whether client is active",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "ISO 8601 timestamp",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *ClientResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	apiClient, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = apiClient
}

// Create creates a new client.
func (r *ClientResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ClientResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the client
	inferaClient, err := r.client.CreateClient(ctx, data.OrganizationID.ValueString(), client.CreateClientRequest{
		Name:    data.Name.ValueString(),
		VaultID: data.VaultID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create client, got error: %s", err))
		return
	}

	// Map response to model
	data.ID = types.StringValue(inferaClient.ID.String())
	data.IsActive = types.BoolValue(inferaClient.IsActive)
	data.CreatedAt = types.StringValue(inferaClient.CreatedAt)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the current client state.
func (r *ClientResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ClientResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the client
	inferaClient, err := r.client.GetClient(ctx, data.OrganizationID.ValueString(), data.ID.ValueString())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read client, got error: %s", err))
		return
	}

	// Map response to model
	data.VaultID = types.StringValue(inferaClient.VaultID.String())
	data.Name = types.StringValue(inferaClient.Name)
	data.IsActive = types.BoolValue(inferaClient.IsActive)
	data.CreatedAt = types.StringValue(inferaClient.CreatedAt)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates an existing client.
func (r *ClientResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ClientResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the client
	inferaClient, err := r.client.UpdateClient(ctx, data.OrganizationID.ValueString(), data.ID.ValueString(), client.UpdateClientRequest{
		Name:    data.Name.ValueString(),
		VaultID: data.VaultID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update client, got error: %s", err))
		return
	}

	// Map response to model
	data.VaultID = types.StringValue(inferaClient.VaultID.String())
	data.Name = types.StringValue(inferaClient.Name)
	data.IsActive = types.BoolValue(inferaClient.IsActive)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes a client.
func (r *ClientResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ClientResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the client
	err := r.client.DeleteClient(ctx, data.OrganizationID.ValueString(), data.ID.ValueString())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			// Resource already deleted, no error
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete client, got error: %s", err))
		return
	}
}

// ImportState imports a client using the format: <org_id>/<client_id>
func (r *ClientResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in format: <org_id>/<client_id>, got: %s", req.ID),
		)
		return
	}

	orgID := parts[0]
	clientID := parts[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), orgID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), clientID)...)
}
