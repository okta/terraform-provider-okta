package okta4ai_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaMcpServersAuthorizationServer_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaMcpServersAuthorizationServer, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	dataSourceName := "data.okta_mcp_servers_authorization_server.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "ausu6ayw5WeizA26Vao5"),
					resource.TestCheckResourceAttr(dataSourceName, "mcp_server_id", "emsu6mjxY5YyMHF0zao5"),
					resource.TestCheckResourceAttrSet(dataSourceName, "created"),
					resource.TestCheckResourceAttrSet(dataSourceName, "last_updated"),
					resource.TestCheckResourceAttrSet(dataSourceName, "authorization_server.status"),
				),
			},
		},
	})
}
