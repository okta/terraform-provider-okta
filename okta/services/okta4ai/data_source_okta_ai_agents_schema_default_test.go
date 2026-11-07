package okta4ai_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaAiAgentsSchemaDefault_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaAIAgentsSchemaDefault, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	dataSourceName := fmt.Sprintf("data.%s.test", resources.OktaAIAgentsSchemaDefault)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					// The org has exactly one AI agent schema; it is returned as a one-item list.
					resource.TestCheckResourceAttr(dataSourceName, "id", "okta_ai_agents_schema_default_list"),
					resource.TestCheckResourceAttr(dataSourceName, "items.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "items.0.id", "ai_agents_schema_default"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.type"),
				),
			},
		},
	})
}
