package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaAuthServerKeys_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSAuthServerKeys, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_auth_server_keys.example", "alg", "RS256"),
					resource.TestCheckResourceAttr("data.okta_auth_server_keys.example", "e", "AQAB"),
					resource.TestCheckResourceAttr("data.okta_auth_server_keys.example", "kid", "abcdefghijk0123456789"),
					resource.TestCheckResourceAttr("data.okta_auth_server_keys.example", "n", "123abc456def"),
				),
			},
		},
	})
}
