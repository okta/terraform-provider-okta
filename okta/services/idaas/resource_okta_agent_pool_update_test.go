package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaAgentPoolUpdate_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAgentPoolUpdate, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.example", resources.OktaIDaaSAgentPoolUpdate)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "schedule_test"),
					resource.TestCheckResourceAttr(resourceName, "agent_type", "AD"),
					resource.TestCheckResourceAttr(resourceName, "notify_admins", "true"),
					resource.TestCheckResourceAttr(resourceName, "pool_id", "0oaspf3cfatE1nDO31d7"),
					resource.TestCheckResourceAttr(resourceName, "schedule.cron", "0 3 * * WED"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "schedule.cron", "0 3 * * MON"),
				),
			},
		},
	})
}
