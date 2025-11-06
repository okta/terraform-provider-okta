package governance_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccRequestSettingResource_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaGovernanceRequestSettingResource, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaGovernanceRequestSettingResource)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				ImportState:        true,
				ResourceName:       "okta_request_setting_resource.test",
				ImportStateId:      "0oaoum6j3cElINe1z1d7",
				ImportStatePersist: true,
				Config:             config,
				PlanOnly:           true,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "risk_settings.default_setting.approval_sequence_id", "68920b41386747a673869356"),
					resource.TestCheckResourceAttr(resourceName, "risk_settings.default_setting.request_submission_type", "ALLOWED_WITH_OVERRIDES"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "risk_settings.default_setting.request_submission_type", "RESTRICTED"),
				),
			},
		},
	})

}
