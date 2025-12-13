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

func TestAccTeamResource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test")
	teamName := acctest.RandomWithPrefix("tf-team")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTeamResourceConfig(rName, teamName, "Test team description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("inferadb_team.test", "name", teamName),
					resource.TestCheckResourceAttr("inferadb_team.test", "description", "Test team description"),
					resource.TestCheckResourceAttrSet("inferadb_team.test", "id"),
					resource.TestCheckResourceAttrSet("inferadb_team.test", "organization_id"),
					resource.TestCheckResourceAttrSet("inferadb_team.test", "created_at"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "inferadb_team.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["inferadb_team.test"]
					if !ok {
						return "", fmt.Errorf("resource not found: inferadb_team.test")
					}
					return fmt.Sprintf("%s/%s", rs.Primary.Attributes["organization_id"], rs.Primary.ID), nil
				},
			},
			// Update testing - change name and description
			{
				Config: testAccTeamResourceConfig(rName, teamName+"-updated", "Updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("inferadb_team.test", "name", teamName+"-updated"),
					resource.TestCheckResourceAttr("inferadb_team.test", "description", "Updated description"),
				),
			},
			// Delete testing is automatic
		},
	})
}

func testAccTeamResourceConfig(orgName, teamName, description string) string {
	return fmt.Sprintf(`
resource "inferadb_organization" "test" {
  name = %[1]q
  tier = "dev"
}

resource "inferadb_team" "test" {
  organization_id = inferadb_organization.test.id
  name            = %[2]q
  description     = %[3]q
}
`, orgName, teamName, description)
}
