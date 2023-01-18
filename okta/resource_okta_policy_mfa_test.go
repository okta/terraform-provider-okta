package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Note: at least one factor (e.g. `okta_otp`) should be enabled before running this test.
func TestAccOktaMfaPolicy_crud(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(policyMfa)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", policyMfa)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createPolicyCheckDestroy(policyMfa),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test MFA Policy"),
					resource.TestCheckResourceAttr(resourceName, "okta_email.enroll", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "google_otp.enroll", "REQUIRED"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)+"_new"),
					resource.TestCheckResourceAttr(resourceName, "status", statusInactive),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test MFA Policy Updated"),
					resource.TestCheckResourceAttr(resourceName, "okta_email.enroll", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "okta_password.enroll", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceName, "google_otp.enroll", "OPTIONAL"),
				),
			},
		},
	})
}

// TestAccOktaMfaPolicy_PR_1210 deals with testing
// https://github.com/okta/terraform-provider-okta/pull/1210
func TestAccOktaMfaPolicy_PR_1210(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(policyMfa)
	config := `
data "okta_group" "all" {
  name = "Everyone"
}
resource "okta_policy_mfa" "test" {
  name        = "testAcc_replace_with_uuid"
  status = "ACTIVE"
  description = "Terraform Acceptance Test MFA Policy"
  priority = 1
  is_oie  = true

  okta_password = {
    enroll = "REQUIRED"
  }

  okta_email = {
    enroll = "NOT_ALLOWED"
  }

  fido_webauthn = {
    enroll = "REQUIRED"
  }

  groups_included = [data.okta_group.all.id]
}
	`
	resourceName := fmt.Sprintf("%s.test", policyMfa)

	resource.Test(t, resource.TestCase{
		PreCheck:          testOIEOnlyAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createPolicyCheckDestroy(policyMfa),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config, ri),
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test MFA Policy"),
					resource.TestCheckResourceAttr(resourceName, "okta_password.enroll", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "okta_email.enroll", "NOT_ALLOWED"),
					resource.TestCheckResourceAttr(resourceName, "fido_webauthn.enroll", "REQUIRED"),
				),
			},
		},
	})
}
