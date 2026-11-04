package okta4ai_test

import (
	"github.com/okta/terraform-provider-okta/okta/resources"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

func TestAccDataSourceOktaMcpServer_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaMcpServer, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	dataSourceName := "data.okta_mcp_server.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "orn"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_url"),
					resource.TestCheckResourceAttrSet(dataSourceName, "created"),
					resource.TestCheckResourceAttrSet(dataSourceName, "last_updated"),
				),
			},
		},
	})
}
