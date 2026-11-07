package okta4ai_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceAiAgentProvidersStagingStat_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaAIAgentProvidersStagingStat, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	dataSourceName := "data." + resources.OktaAIAgentProvidersStagingStat + ".test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "provider_id"),
					resource.TestCheckResourceAttr(dataSourceName, "items.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "items.0.id", "ai_agent_providers_staging_stat"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.duplicate_count"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.ignored_count"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.new_count"),
				),
			},
		},
	})
}
