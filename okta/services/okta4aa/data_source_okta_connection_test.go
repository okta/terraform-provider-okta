package okta4aa_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaConnection_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaConnection, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	dataSourceName := fmt.Sprintf("data.%s.test", resources.OktaConnection)
	resourceName := fmt.Sprintf("%s.test", resources.OktaConnection)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "orn", resourceName, "orn"),
					resource.TestCheckResourceAttr(dataSourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(dataSourceName, "actor.wlp_id", resourceName, "actor.wlp_id"),
					resource.TestCheckResourceAttr(dataSourceName, "actor.wlp_type", "AI_AGENT"),
					resource.TestCheckResourceAttr(dataSourceName, "authorization_server.orn", "orn:okta:idp:00otxk5i5T0Rk59Z3ao5:authorization_servers:ausu6ayw5WeizA26Vao5"),
					resource.TestCheckResourceAttrSet(dataSourceName, "authorization_server.issuer_url"),
					// Union-variant enums: read through a per-variant type switch.
					resource.TestCheckResourceAttr(dataSourceName, "protocol.protocol_type", "IDENTITY_ASSERTION"),
					resource.TestCheckResourceAttr(dataSourceName, "protocol.scope_condition", "INCLUDE_ONLY"),
					resource.TestCheckResourceAttr(dataSourceName, "protocol.scopes.#", "3"),
					resource.TestCheckResourceAttr(dataSourceName, "resource.resource_type", "MCP_SERVER"),
					resource.TestCheckResourceAttr(dataSourceName, "resource.orn", "orn:okta:directory:00otxk5i5T0Rk59Z3ao5:resource-servers:mcp:emsv9wtzWEC6F5C0mao5"),
					resource.TestCheckResourceAttr(dataSourceName, "resource.display_name", "Microsoft Teams"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_identifier"),
				),
			},
		},
	})
}
