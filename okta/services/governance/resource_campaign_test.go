package governance_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

// Warning: Users need to enable governance manually from the console for applications, since it can't be enabled via API due to API not being public.
func TestAccCampaignResource_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaGovernanceCampaign, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaGovernanceCampaign)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Monthly access review of sales team"),
					resource.TestCheckResourceAttr(resourceName, "campaign_type", "RESOURCE"),
					resource.TestCheckResourceAttr(resourceName, "schedule_settings.type", "ONE_OFF"),
					resource.TestCheckResourceAttr(resourceName, "principal_scope_settings.type", "USERS"),
				),
			},
		},
	})
}
