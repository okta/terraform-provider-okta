package okta4ai_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaAiAgentsCredentialsSecret_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaAIAgentsCredentialsSecret, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaAIAgentsCredentialsSecret)

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
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet(resourceName, "secret_hash"),
					// The API returns the plaintext secret only on create.
					resource.TestCheckResourceAttrSet(resourceName, "client_secret"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: importStateIdAgentChild(resourceName),
				ImportStateVerify: true,
				// The plaintext secret is returned only at creation, so a read after import cannot restore it.
				ImportStateVerifyIgnore: []string{"client_secret"},
			},
		},
	})
}
