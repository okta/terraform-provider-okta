package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

func TestAccDataSourceOktaPoliciesRuleAccessPolicy_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", "okta_policies_rule_access_policy", t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_policies_rule_access_policy.test", "id"),
					resource.TestCheckResourceAttrSet("data.okta_policies_rule_access_policy.test", "name"),
					resource.TestCheckResourceAttrSet("data.okta_policies_rule_access_policy.test", "status"),
					resource.TestCheckResourceAttrSet("data.okta_policies_rule_access_policy.test", "priority"),
					resource.TestCheckResourceAttr("data.okta_policies_rule_access_policy.test", "status", "ACTIVE"),
				),
			},
		},
	})
}
