package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaAppToken_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSAppToken, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_app_token.example", "id"),
					resource.TestCheckResourceAttr("data.okta_app_token.example", "status", "ACTIVE"),
					resource.TestCheckResourceAttr("data.okta_app_token.example", "scopes.#", "2"),
				),
			},
		},
	})
}
