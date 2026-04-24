package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaAgentPool_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSAgentPoolUpdate, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_agent_pool_update.example", "id"),
					resource.TestCheckResourceAttrSet("data.okta_agent_pool_update.example", "pool_id"),
					resource.TestCheckResourceAttr("data.okta_agent_pool_update.example", "name", "update_schedule_test"),
					resource.TestCheckResourceAttr("data.okta_agent_pool_update.example", "agent_type", "AD"),
					resource.TestCheckResourceAttr("data.okta_agent_pool_update.example", "enabled", "false"),
					resource.TestCheckResourceAttr("data.okta_agent_pool_update.example", "notify_admin", "true"),
					resource.TestCheckResourceAttr("data.okta_agent_pool_update.example", "status", "Scheduled"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"data.okta_agent_pool_update.example",
						"agents.*",
						map[string]string{
							"operational_status": "OPERATIONAL",
							"type":               "AD",
						},
					),
				),
			},
		},
	})
}
