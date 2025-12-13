// Copyright 2025 InferaDB
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccClientResource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test")
	clientName := acctest.RandomWithPrefix("tf-client")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccClientResourceConfig(rName, clientName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("inferadb_client.test", "name", clientName),
					resource.TestCheckResourceAttrSet("inferadb_client.test", "id"),
					resource.TestCheckResourceAttrSet("inferadb_client.test", "organization_id"),
					resource.TestCheckResourceAttrSet("inferadb_client.test", "vault_id"),
					resource.TestCheckResourceAttr("inferadb_client.test", "is_active", "true"),
					resource.TestCheckResourceAttrSet("inferadb_client.test", "created_at"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "inferadb_client.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["inferadb_client.test"]
					if !ok {
						return "", fmt.Errorf("resource not found: inferadb_client.test")
					}
					return fmt.Sprintf("%s/%s", rs.Primary.Attributes["organization_id"], rs.Primary.ID), nil
				},
			},
			// Update testing - change name
			{
				Config: testAccClientResourceConfig(rName, clientName+"-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("inferadb_client.test", "name", clientName+"-updated"),
				),
			},
			// Delete testing is automatic
		},
	})
}

func testAccClientResourceConfig(orgName, clientName string) string {
	return fmt.Sprintf(`
resource "inferadb_organization" "test" {
  name = %[1]q
  tier = "dev"
}

resource "inferadb_vault" "test" {
  organization_id = inferadb_organization.test.id
  name            = "test-vault"
}

resource "inferadb_client" "test" {
  organization_id = inferadb_organization.test.id
  vault_id        = inferadb_vault.test.id
  name            = %[2]q
}
`, orgName, clientName)
}
