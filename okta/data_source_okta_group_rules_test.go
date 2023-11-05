package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOktaGroupRules_read(t *testing.T) {
	mgr := newFixtureManager("datasources", groupRules, t.Name())
	rulesCreate := groupAndRules
	rulesRead := fmt.Sprintf("%s%s", groupAndRule, groupRuleDataSources)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(rulesCreate),
			},
			{
				Config: mgr.ConfigReplace(rulesRead),
				Check: resource.ComposeTestCheckFunc(
					// hack for eventual consistency on the group rule creation on Okta API side
					sleepInSecondsForTest(2),

					resource.TestCheckResourceAttr("data.okta_group_rules.test_by_exact_match", "rules.#", "1"),
					resource.TestCheckResourceAttrSet("data.okta_group_rules.test_by_exact_match", "rules.#.id"),
					resource.TestCheckResourceAttr("data.okta_group_rules.test_by_exact_match", "name", fmt.Sprintf("testRule_%s_one", buildResourceName(mgr.Seed))),
					resource.TestCheckResourceAttr("data.okta_group_rules.test_by_exact_match", "status", "ACTIVE"),

					resource.TestCheckResourceAttr("data.okta_group_rules.test_by_prefix", "rules.#", "2"),

					resource.TestCheckResourceAttr("data.okta_group_rules.test_by_no_match", "rules.#", "0"),
				),
			},
		},
	})
}

const groupAndRules = `
resource "okta_group" "test" {
  name = "testAcc_replace_with_uuid"
}

resource "okta_group_rule" "test1" {
  name              = "testRule_testAcc_replace_with_uuid_one"
  status            = "ACTIVE"
  group_assignments = [okta_group.test.id]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "String.startsWith(user.firstName,\"andy\")"
}

resource "okta_group_rule" "test2" {
  name              = "testRule_testAcc_replace_with_uuid_two"
  status            = "ACTIVE"
  group_assignments = [okta_group.test.id]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "String.startsWith(user.firstName,\"andy\")"
  depends_on        = [okta_group_rule.test1]
}
`

const groupRulesDataSources = `

data "okta_group_rules" "test_by_exact_match" {
  name_prefix = "testRule_testAcc_replace_with_uuid_one"
}

data "okta_group_rules" "test_by_prefix" {
  name_prefix = "testRule_testAcc_replace_with_uuid_"
}

data "okta_group_rules" "test_by_no_match" {
  name_prefix = "invalidRule_replace_with_uuid"
}
`
