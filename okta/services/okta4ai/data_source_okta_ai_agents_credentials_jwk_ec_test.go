package okta4ai_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaAiAgentsCredentialsJwkEc_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaAIAgentsCredentialsJwkEc, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	dataSourceName := fmt.Sprintf("data.%s.test", resources.OktaAIAgentsCredentialsJwkEc)
	resourceName := fmt.Sprintf("%s.test", resources.OktaAIAgentsCredentialsJwkEc)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "agent_id", "okta_ai_agent.test", "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "crv", "P-256"),
					resource.TestCheckResourceAttr(dataSourceName, "x", "WKn-ZIGevcwGIyyrzFoZNBdaq9_TsqzGl96oc0CWuis"),
					resource.TestCheckResourceAttr(dataSourceName, "y", "y77t-RvAHRKTsSGdIYUfweuOvwrvDD-Q3Hv5J0fSKQE"),
					resource.TestCheckResourceAttrPair(dataSourceName, "kid", resourceName, "kid"),
					resource.TestCheckResourceAttrSet(dataSourceName, "status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "created"),
				),
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
