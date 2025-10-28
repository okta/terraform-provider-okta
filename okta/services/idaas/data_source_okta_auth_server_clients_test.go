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
					resource.TestCheckResourceAttrSet("data.okta_auth_server_clients.test", "id"),
					resource.TestCheckResourceAttrSet("data.okta_auth_server_clients.test", "auth_server_id"),
					resource.TestCheckResourceAttrSet("data.okta_auth_server_clients.test", "client_id"),
					resource.TestCheckResourceAttrSet("data.okta_auth_server_clients.test", "client_name"),
				),
			},
		},
	})
}
