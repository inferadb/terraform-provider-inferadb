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

func TestAccVaultResource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test")
	vaultName := acctest.RandomWithPrefix("tf-vault")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccVaultResourceConfig(rName, vaultName, "Test vault description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("inferadb_vault.test", "name", vaultName),
					resource.TestCheckResourceAttr("inferadb_vault.test", "description", "Test vault description"),
					resource.TestCheckResourceAttrSet("inferadb_vault.test", "id"),
					resource.TestCheckResourceAttrSet("inferadb_vault.test", "organization_id"),
					resource.TestCheckResourceAttrSet("inferadb_vault.test", "created_at"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "inferadb_vault.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["inferadb_vault.test"]
					if !ok {
						return "", fmt.Errorf("resource not found: inferadb_vault.test")
					}
					return fmt.Sprintf("%s/%s", rs.Primary.Attributes["organization_id"], rs.Primary.ID), nil
				},
			},
			// Update testing - change name and description
			{
				Config: testAccVaultResourceConfig(rName, vaultName+"-updated", "Updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("inferadb_vault.test", "name", vaultName+"-updated"),
					resource.TestCheckResourceAttr("inferadb_vault.test", "description", "Updated description"),
				),
			},
			// Delete testing is automatic
		},
	})
}

func testAccVaultResourceConfig(orgName, vaultName, description string) string {
	return fmt.Sprintf(`
resource "inferadb_organization" "test" {
  name = %[1]q
  tier = "dev"
}

resource "inferadb_vault" "test" {
  organization_id = inferadb_organization.test.id
  name            = %[2]q
  description     = %[3]q
}
`, orgName, vaultName, description)
}
