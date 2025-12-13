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
var _ resource.Resource = &ClientCertificateResource{}
var _ resource.ResourceWithImportState = &ClientCertificateResource{}

// NewClientCertificateResource creates a new client certificate resource.
func NewClientCertificateResource() resource.Resource {
	return &ClientCertificateResource{}
}

// ClientCertificateResource defines the resource implementation.
type ClientCertificateResource struct {
	client *client.Client
}

// ClientCertificateResourceModel describes the resource data model.
type ClientCertificateResourceModel struct {
	ID             types.String `tfsdk:"id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	ClientID       types.String `tfsdk:"client_id"`
	Name           types.String `tfsdk:"name"`
	KID            types.String `tfsdk:"kid"`
	PublicKeyPEM   types.String `tfsdk:"public_key_pem"`
	PrivateKeyPEM  types.String `tfsdk:"private_key_pem"`
	IsActive       types.Bool   `tfsdk:"is_active"`
	RevokedAt      types.String `tfsdk:"revoked_at"`
	CreatedAt      types.String `tfsdk:"created_at"`
}

// Metadata returns the resource type name.
func (r *ClientCertificateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_client_certificate"
}

// Schema defines the schema for the resource.
func (r *ClientCertificateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Manages an InferaDB client certificate.

Client certificates are Ed25519 key pairs used for JWT-based authentication.
Each certificate is associated with a client (backend service identity) and can be used
to generate JWTs for API authentication.

**IMPORTANT**: The private key is only returned during creation and cannot be retrieved later.
Store it securely immediately after creation.

## Example Usage

` + "```hcl" + `
resource "inferadb_client_certificate" "backend_cert" {
  organization_id = inferadb_organization.example.id
  client_id       = inferadb_client.backend.id
  name            = "Production Backend Certificate"
}

# Save the private key to a secure location
output "backend_private_key" {
  value     = inferadb_client_certificate.backend_cert.private_key_pem
  sensitive = true
}
` + "```",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique Snowflake ID of the certificate.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "ID of the organization that owns the client.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "ID of the client this certificate belongs to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Human-readable name for the certificate (e.g., 'Production Backend Cert').",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"kid": schema.StringAttribute{
				MarkdownDescription: "Key ID (kid) used in JWT headers for signature verification.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"public_key_pem": schema.StringAttribute{
				MarkdownDescription: "Ed25519 public key in PEM format.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"private_key_pem": schema.StringAttribute{
				MarkdownDescription: "Ed25519 private key in PEM format. **CRITICAL**: This is only returned on creation and cannot be retrieved later. Store securely.",
				Computed:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_active": schema.BoolAttribute{
				MarkdownDescription: "Whether the certificate is active and can be used for authentication.",
				Computed:            true,
			},
			"revoked_at": schema.StringAttribute{
				MarkdownDescription: "ISO 8601 timestamp when the certificate was revoked (null if active).",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "ISO 8601 timestamp when the certificate was created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure configures the resource with the provider client.
func (r *ClientCertificateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a new client certificate.
func (r *ClientCertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ClientCertificateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the certificate
	cert, err := r.client.CreateCertificate(ctx, data.OrganizationID.ValueString(), data.ClientID.ValueString(), client.CreateCertificateRequest{
		Name: data.Name.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Client Certificate",
			fmt.Sprintf("Could not create client certificate: %s", err.Error()),
		)
		return
	}

	// Map response to model
	data.ID = types.StringValue(cert.ID)
	data.KID = types.StringValue(cert.KID)
	data.PublicKeyPEM = types.StringValue(cert.PublicKeyPEM)
	data.IsActive = types.BoolValue(cert.IsActive)
	data.CreatedAt = types.StringValue(cert.CreatedAt)

	// CRITICAL: Private key is only returned on creation
	if cert.PrivateKeyPEM != "" {
		data.PrivateKeyPEM = types.StringValue(cert.PrivateKeyPEM)
	} else {
		data.PrivateKeyPEM = types.StringNull()
	}

	// Handle optional fields
	if cert.RevokedAt != nil {
		data.RevokedAt = types.StringValue(*cert.RevokedAt)
	} else {
		data.RevokedAt = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the resource state.
func (r *ClientCertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ClientCertificateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current certificate state
	cert, err := r.client.GetCertificate(ctx, data.OrganizationID.ValueString(), data.ClientID.ValueString(), data.ID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			// Certificate was deleted outside Terraform
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Client Certificate",
			fmt.Sprintf("Could not read client certificate %s: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	// Update state with current values
	data.KID = types.StringValue(cert.KID)
	data.PublicKeyPEM = types.StringValue(cert.PublicKeyPEM)
	data.IsActive = types.BoolValue(cert.IsActive)
	data.CreatedAt = types.StringValue(cert.CreatedAt)

	// Handle optional fields
	if cert.RevokedAt != nil {
		data.RevokedAt = types.StringValue(*cert.RevokedAt)
	} else {
		data.RevokedAt = types.StringNull()
	}

	// IMPORTANT: Private key is NOT returned on reads, preserve state value
	// The private key in state will remain from creation

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource (no-op for certificates as all fields require replacement).
func (r *ClientCertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// All certificate attributes require replacement, so this should never be called
	// This is enforced by RequiresReplace plan modifiers in the schema
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Client certificates cannot be updated. All changes require replacement.",
	)
}

// Delete deletes the resource.
func (r *ClientCertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ClientCertificateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the certificate
	err := r.client.DeleteCertificate(ctx, data.OrganizationID.ValueString(), data.ClientID.ValueString(), data.ID.ValueString())
	if err != nil {
		// Ignore 404 errors as the resource may have been deleted outside Terraform
		if !strings.Contains(err.Error(), "404") && !strings.Contains(err.Error(), "not found") {
			resp.Diagnostics.AddError(
				"Error Deleting Client Certificate",
				fmt.Sprintf("Could not delete client certificate %s: %s", data.ID.ValueString(), err.Error()),
			)
			return
		}
	}
}

// ImportState imports an existing resource into Terraform state.
func (r *ClientCertificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expected format: <org_id>/<client_id>/<cert_id>
	parts := strings.Split(req.ID, "/")
	if len(parts) != 3 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID format: <org_id>/<client_id>/<cert_id>, got: %s", req.ID),
		)
		return
	}

	orgID := parts[0]
	clientID := parts[1]
	certID := parts[2]

	// Set the IDs in state
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), orgID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("client_id"), clientID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), certID)...)

	// Note: After import, the Read method will be called automatically to populate the rest of the state.
	// The private_key_pem will be null since it cannot be retrieved after creation.
}
