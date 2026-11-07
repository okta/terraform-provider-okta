package okta4ai_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaAiAgentsCredentialsJwksUri_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaAIAgentsCredentialsJwksUri, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaAIAgentsCredentialsJwksUri)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "agent_id", "okta_ai_agent.test", "id"),
					// One per agent: the ID is the agent ID.
					resource.TestCheckResourceAttrPair(resourceName, "id", "okta_ai_agent.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "jwks_uri", "https://example.com/.well-known/jwks.json"),
				),
			},
			{
				// Import takes just the agent ID.
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return rs.Primary.Attributes["agent_id"], nil
				},
				ImportStateVerify: true,
			},
		},
	})
}
