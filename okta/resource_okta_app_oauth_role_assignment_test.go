package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaAppOAuthRoleAssignment_basic(t *testing.T) {
	mgr := newFixtureManager("okta_app_oauth_role_assignment", t.Name())

	oktaResourceTest(t, resource.TestCase{
		PreCheck:                 testAccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
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
	mgr := newFixtureManager("okta_app_oauth_role_assignment", t.Name())

	oktaResourceTest(t, resource.TestCase{
		PreCheck:                 testAccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
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
