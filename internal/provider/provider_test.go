// Copyright 2025 InferaDB
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"inferadb": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// Check that required environment variables are set
	if os.Getenv("INFERADB_SESSION_TOKEN") == "" {
		t.Fatal("INFERADB_SESSION_TOKEN must be set for acceptance tests")
	}
	if os.Getenv("INFERADB_ENDPOINT") == "" {
		t.Fatal("INFERADB_ENDPOINT must be set for acceptance tests")
	}
}
