package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOktaGroupRule_read(t *testing.T) {
	mgr := newFixtureManager(groupRule, t.Name())
	step1config := groupAndRule
	step2config := fmt.Sprintf("%s%s", groupAndRule, groupRuleDataSources)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(step1config),
			},
			{
				Config: mgr.ConfigReplace(step2config),
				Check: resource.ComposeTestCheckFunc(
					// hack for eventual consistency on the group rule creation on Okta API side
					sleepInSecondsForTest(2),

					resource.TestCheckResourceAttrSet("data.okta_group_rule.test_by_id", "id"),
					resource.TestCheckResourceAttr("data.okta_group_rule.test_by_id", "name", fmt.Sprintf("testAccTwo_%d", mgr.Seed)),
					resource.TestCheckResourceAttr("data.okta_group_rule.test_by_id", "status", "ACTIVE"),

					resource.TestCheckResourceAttrSet("data.okta_group_rule.test_by_name", "id"),
					resource.TestCheckResourceAttr("data.okta_group_rule.test_by_name", "name", fmt.Sprintf("testAccTwo_%d", mgr.Seed)),
					resource.TestCheckResourceAttr("data.okta_group_rule.test_by_name", "status", "ACTIVE"),
				),
			},
		},
	})
}

const groupAndRule = `
resource "okta_group" "test" {
  name = "testAcc_replace_with_uuid"
}

resource "okta_group_rule" "test1" {
  name              = "testAccOne_replace_with_uuid"
  status            = "ACTIVE"
  group_assignments = [okta_group.test.id]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "String.startsWith(user.firstName,\"andy\")"
}

resource "okta_group_rule" "test2" {
  name              = "testAccTwo_replace_with_uuid"
  status            = "ACTIVE"
  group_assignments = [okta_group.test.id]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "String.startsWith(user.firstName,\"andy\")"
  depends_on        = [okta_group_rule.test1]
}
`

const groupRuleDataSources = `
data "okta_group_rule" "test_by_id" {
  id          = okta_group_rule.test2.id
}

data "okta_group_rule" "test_by_name" {
  name          = "testAccTwo_replace_with_uuid"
}
`
