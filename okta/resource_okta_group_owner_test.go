package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaGroupOwner_crud(t *testing.T) {
	mgr := newFixtureManager("resources", "okta_group_owner", t.Name())
	config := mgr.GetFixtures("resource.tf", t)

	oktaResourceTest(
		t, resource.TestCase{
			PreCheck:                 testAccPreCheck(t),
			ErrorCheck:               testAccErrorChecks(t),
			CheckDestroy:             nil,
			ProtoV5ProviderFactories: testAccMergeProvidersFactories,
			Steps: []resource.TestStep{
				{
					Config: config,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("okta_user.test", "first_name", "TestAcc"),
						resource.TestCheckResourceAttr("okta_user.test", "last_name", "Smith"),
						resource.TestCheckResourceAttr("okta_group_owner.test", "type", "USER"),
					),
				},
			},
		})
}
