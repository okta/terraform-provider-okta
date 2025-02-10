package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceOktaAppOAuthRoleAssignment_basic(t *testing.T) {
	mgr := newFixtureManager("resources", appOAuthRoleAssignment, t.Name())

	oktaResourceTest(t, resource.TestCase{
		PreCheck:                 testAccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		CheckDestroy:             nil,
		ProtoV5ProviderFactories: testAccMergeProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: mgr.GetFixtures("basic.tf", t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("okta_app_oauth_role_assignment.test", "type", "HELP_DESK_ADMIN"),
					resource.TestCheckResourceAttr("okta_app_oauth_role_assignment.test", "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet("okta_app_oauth_role_assignment.test", "id"),
					resource.TestCheckResourceAttrSet("okta_app_oauth_role_assignment.test", "client_id"),
					resource.TestCheckResourceAttrSet("okta_app_oauth_role_assignment.test", "label"),
				),
			},
			{
				ResourceName:      "okta_app_oauth_role_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					r, ok := s.RootModule().Resources["okta_app_oauth_role_assignment.test"]
					if !ok {
						return "", fmt.Errorf("Unable to find resource: %s:", "okta_app_oauth_role_assignment.test")
					}
					return fmt.Sprintf("%s/%s", r.Primary.Attributes["client_id"], r.Primary.Attributes["id"]), nil
				},
			},
			{
				Config: mgr.GetFixtures("basic_updated.tf", t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("okta_app_oauth_role_assignment.test", "type", "GROUP_MEMBERSHIP_ADMIN"),
					resource.TestCheckResourceAttr("okta_app_oauth_role_assignment.test", "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet("okta_app_oauth_role_assignment.test", "id"),
					resource.TestCheckResourceAttrSet("okta_app_oauth_role_assignment.test", "client_id"),
					resource.TestCheckResourceAttrSet("okta_app_oauth_role_assignment.test", "label"),
				),
			},
		},
	})
}

func TestAccResourceOktaAppOAuthRoleAssignment_custom(t *testing.T) {
	mgr := newFixtureManager("resources", appOAuthRoleAssignment, t.Name())

	oktaResourceTest(t, resource.TestCase{
		PreCheck:                 testAccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		CheckDestroy:             nil,
		ProtoV5ProviderFactories: testAccMergeProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: mgr.GetFixtures("custom.tf", t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("okta_app_oauth_role_assignment.test", "type", "CUSTOM"),
					resource.TestCheckResourceAttr("okta_app_oauth_role_assignment.test", "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet("okta_app_oauth_role_assignment.test", "id"),
					resource.TestCheckResourceAttrSet("okta_app_oauth_role_assignment.test", "client_id"),
					resource.TestCheckResourceAttrSet("okta_app_oauth_role_assignment.test", "label"),
					resource.TestCheckResourceAttrSet("okta_app_oauth_role_assignment.test", "role"),
					resource.TestCheckResourceAttrSet("okta_app_oauth_role_assignment.test", "resource_set"),
				),
			},
			{
				Config: mgr.GetFixtures("custom_updated.tf", t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("okta_app_oauth_role_assignment.test", "type", "CUSTOM"),
					resource.TestCheckResourceAttr("okta_app_oauth_role_assignment.test", "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet("okta_app_oauth_role_assignment.test", "id"),
					resource.TestCheckResourceAttrSet("okta_app_oauth_role_assignment.test", "client_id"),
					resource.TestCheckResourceAttrSet("okta_app_oauth_role_assignment.test", "label"),
					resource.TestCheckResourceAttrSet("okta_app_oauth_role_assignment.test", "role"),
					resource.TestCheckResourceAttrSet("okta_app_oauth_role_assignment.test", "resource_set"),
				),
			},
		},
	})
}
