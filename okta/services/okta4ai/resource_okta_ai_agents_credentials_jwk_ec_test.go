package okta4ai_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaAiAgentsCredentialsJwkEc_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaAIAgentsCredentialsJwkEc, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaAIAgentsCredentialsJwkEc)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "agent_id", "okta_ai_agent.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "kty", "EC"),
					resource.TestCheckResourceAttr(resourceName, "crv", "P-256"),
					resource.TestCheckResourceAttr(resourceName, "alg", "ES256"),
					resource.TestCheckResourceAttr(resourceName, "use", "sig"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "x", "WKn-ZIGevcwGIyyrzFoZNBdaq9_TsqzGl96oc0CWuis"),
					resource.TestCheckResourceAttr(resourceName, "y", "y77t-RvAHRKTsSGdIYUfweuOvwrvDD-Q3Hv5J0fSKQE"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: importStateIdAgentChild(resourceName),
				ImportStateVerify: true,
			},
			{
				// Skip deleting the key. Okta won't deactivate an agent's last active key and only
				// deletes inactive keys, so the key can't be destroyed on its own. removed.tf forgets
				// it without destroying it; the post-test destroy deletes the agent, which removes it.
				Config: mgr.GetFixtures("removed.tf", t),
			},
		},
	})
}

// importStateIdAgentChild builds the {agent_id}/{id} import ID used by AI agent credential resources.
func importStateIdAgentChild(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["agent_id"], rs.Primary.ID), nil
	}
}
