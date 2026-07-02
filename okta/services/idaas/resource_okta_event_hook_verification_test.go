package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaEventHookVerification_crud(t *testing.T) {
	resourceName := "okta_event_hook_verification.user_assigned"
	mgr := newFixtureManager("resources", resources.OktaIDaaSEventHookVerification, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "verification_status", "VERIFIED"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "verification_status", "VERIFIED"),
				),
				// After apply, ReadContext fetches UNVERIFIED from the API and
				// CustomizeDiff forces the plan to VERIFIED, so a non-empty plan
				// is expected on the post-step refresh.
				ExpectNonEmptyPlan: true,
			},
			{
				// ReadContext fetches the hook and sets verification_status=UNVERIFIED
				// in state (e.g. the hook was reset externally between applies).
				// CustomizeDiff then detects UNVERIFIED and calls
				// d.SetNew("verification_status", "VERIFIED"), so the plan must show ...
				Config:             updatedConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true, // ... at least one change — Terraform will schedule an Update to re-verify.
			},
		},
	})
}
