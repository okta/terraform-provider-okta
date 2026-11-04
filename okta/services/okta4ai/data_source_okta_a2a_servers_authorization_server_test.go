package okta4ai_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaA2aServersAuthorizationServer_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaA2aServersAuthorizationServer, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	dataSourceName := "data.okta_a2a_servers_authorization_server.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "a2a_server_id", "wlpu37lfb7LQOikYDao5"),
					resource.TestCheckResourceAttrSet(dataSourceName, "created"),
					resource.TestCheckResourceAttrSet(dataSourceName, "last_updated"),
					resource.TestCheckResourceAttrSet(dataSourceName, "authorization_server.status"),
				),
			},
		},
	})
}
