package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaAPIToken_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSAPIToken, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_api_token.test", "id"),
					resource.TestCheckResourceAttr("data.okta_api_token.test", "name", "api-token-test-token"),
					resource.TestCheckResourceAttr("data.okta_api_token.test", "network.connection", "ANYWHERE"),
				),
			},
		},
	})
}
