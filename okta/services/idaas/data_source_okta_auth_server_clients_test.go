package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaAuthServerClients_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSAuthServerClients, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_auth_server_clients.test", "id", "oar123456abcdefghijklm"),
					resource.TestCheckResourceAttr("data.okta_auth_server_clients.test", "auth_server_id", "aus123456abcdefghijklm"),
					resource.TestCheckResourceAttr("data.okta_auth_server_clients.test", "client_id", "0oa123456abcdefghijklm"),
					resource.TestCheckResourceAttr("data.okta_auth_server_clients.test", "created", "2025-10-28 17:36:58 +0000 UTC"),
					resource.TestCheckResourceAttr("data.okta_auth_server_clients.test", "expires_at", "2025-11-04 17:36:58 +0000 UTC"),
					resource.TestCheckResourceAttr("data.okta_auth_server_clients.test", "status", "ACTIVE"),
					resource.TestCheckResourceAttr("data.okta_auth_server_clients.test", "user_id", "00uonmvbfznIufFS61d7"),
					resource.TestCheckResourceAttr("data.okta_auth_server_clients.test", "scopes.#", "2"),
				),
			},
		},
	})
}
