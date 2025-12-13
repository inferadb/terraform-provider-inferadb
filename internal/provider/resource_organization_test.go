// Copyright 2025 InferaDB
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganizationResource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccOrganizationResourceConfig(rName, "dev"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("inferadb_organization.test", "name", rName),
					resource.TestCheckResourceAttr("inferadb_organization.test", "tier", "dev"),
					resource.TestCheckResourceAttrSet("inferadb_organization.test", "id"),
					resource.TestCheckResourceAttrSet("inferadb_organization.test", "created_at"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "inferadb_organization.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update testing - change name
			{
				Config: testAccOrganizationResourceConfig(rName+"-updated", "dev"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("inferadb_organization.test", "name", rName+"-updated"),
					resource.TestCheckResourceAttr("inferadb_organization.test", "tier", "dev"),
				),
			},
			// Delete testing is automatic
		},
	})
}

func testAccOrganizationResourceConfig(name, tier string) string {
	return fmt.Sprintf(`
resource "inferadb_organization" "test" {
  name = %[1]q
  tier = %[2]q
}
`, name, tier)
}
