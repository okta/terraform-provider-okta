package okta4ai_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaConnection_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaConnection, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaConnection)
	agentName := fmt.Sprintf("%s.test", resources.OktaAIAgent)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "orn"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(resourceName, "actor.wlp_id", agentName, "id"),
					resource.TestCheckResourceAttr(resourceName, "actor.wlp_type", "AI_AGENT"),
					resource.TestCheckResourceAttrSet(resourceName, "authorization_server.issuer_url"),
					resource.TestCheckResourceAttr(resourceName, "protocol.protocol_type", "IDENTITY_ASSERTION"),
					resource.TestCheckResourceAttr(resourceName, "protocol.scope_condition", "INCLUDE_ONLY"),
					resource.TestCheckResourceAttr(resourceName, "protocol.scopes.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "resource.resource_type", "MCP_SERVER"),
					resource.TestCheckResourceAttr(resourceName, "resource.orn", "orn:okta:directory:00otxk5i5T0Rk59Z3ao5:resource-servers:mcp:emsv9wtzWEC6F5C0mao5"),
					resource.TestCheckResourceAttr(resourceName, "resource.display_name", "Microsoft Teams"),
				),
			},
			{
				// Changing protocol scopes updates the connection in place.
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "protocol.scopes.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "protocol.scopes.1", "profile"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
