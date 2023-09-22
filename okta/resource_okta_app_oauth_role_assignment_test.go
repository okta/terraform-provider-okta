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
			},
			{
				Config: mgr.GetFixtures("basic_updated.tf", t),
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
			},
			{
				Config: mgr.GetFixtures("custom_updated.tf", t),
			},
		},
	})
}
