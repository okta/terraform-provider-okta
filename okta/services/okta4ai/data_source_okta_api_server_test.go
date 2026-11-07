package okta4ai_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaApiServer_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaApiServer, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	dataSourceName := fmt.Sprintf("data.%s.test", resources.OktaApiServer)
	resourceName := fmt.Sprintf("%s.test", resources.OktaApiServer)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "resource_url", resourceName, "resource_url"),
					resource.TestCheckResourceAttrPair(dataSourceName, "orn", resourceName, "orn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "status", resourceName, "status"),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.display_name", fmt.Sprintf("testAcc-api-server-ds-%d", mgr.Seed)),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.description", "Terraform acceptance test API server for the okta_api_server data source"),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.scopes_supported.#", "2"),
					resource.TestCheckResourceAttrSet(dataSourceName, "created"),
					resource.TestCheckResourceAttrSet(dataSourceName, "last_updated"),
				),
			},
		},
	})
}
