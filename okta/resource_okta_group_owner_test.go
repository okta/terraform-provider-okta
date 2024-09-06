package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaGroupOwner_crud(t *testing.T) {

	oktaResourceTest(
		t, resource.TestCase{
			PreCheck:                 testAccPreCheck(t),
			ErrorCheck:               testAccErrorChecks(t),
			ProtoV5ProviderFactories: testAccMergeProvidersFactories,
			Steps: []resource.TestStep{
				{
					Config: `resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc-replace_with_uuid@example.com"
  email      = "testAcc-replace_with_uuid@example.com"
}

resource "okta_group" "test" {
  name = "testAcc_replace_with_uuid"
}

resource "okta_group_owner" "test" {
  group_id                  = okta_group.test.id
  id_of_group_owner         = okta_user.test.id
  type                      = "USER"
}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("okta_user.test", "first_name", "TestAcc"),
						resource.TestCheckResourceAttr("okta_user.test", "last_name", "Smith"),
						resource.TestCheckResourceAttr("okta_group_owner.test", "type", "USER"),
					),
				},
			},
		})
}
