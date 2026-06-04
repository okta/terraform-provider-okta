package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

func TestAccDataSourceOktaAppSignOnPolicyRule_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", "okta_app_sign_on_policy_rule", t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_app_sign_on_policy_rule.test", "id"),
					resource.TestCheckResourceAttrSet("data.okta_app_sign_on_policy_rule.test", "name"),
					resource.TestCheckResourceAttrSet("data.okta_app_sign_on_policy_rule.test", "status"),
					resource.TestCheckResourceAttrSet("data.okta_app_sign_on_policy_rule.test", "priority"),
					resource.TestCheckResourceAttr("data.okta_app_sign_on_policy_rule.test", "status", "ACTIVE"),
				),
			},
		},
	})
}
