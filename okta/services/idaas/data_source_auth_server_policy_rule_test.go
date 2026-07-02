package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

func TestAccDataSourceAuthorizationServersPoliciesRule_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", "okta_authorization_servers_policies_rule", t.Name())
	datasourceConfig := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: datasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_authorization_servers_policies_rule.test", "id"),
					resource.TestCheckResourceAttrSet("data.okta_authorization_servers_policies_rule.test", "name"),
					resource.TestCheckResourceAttrSet("data.okta_authorization_servers_policies_rule.test", "status"),
					resource.TestCheckResourceAttrSet("data.okta_authorization_servers_policies_rule.test", "priority"),
					resource.TestCheckResourceAttr("data.okta_authorization_servers_policies_rule.test", "name", "test"),
					resource.TestCheckResourceAttr("data.okta_authorization_servers_policies_rule.test", "status", "ACTIVE"),
				),
			},
		},
	})
}
