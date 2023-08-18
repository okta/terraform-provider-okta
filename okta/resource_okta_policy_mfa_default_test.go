package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaDefaultMFAPolicy(t *testing.T) {
	mgr := newFixtureManager(policyMfaDefault, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", policyMfaDefault)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "okta_password.enroll", "REQUIRED"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "okta_password.enroll", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "google_otp.enroll", "OPTIONAL"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "okta_password.enroll", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
				),
			},
		},
	})
}

// TestAccResourceOktaMfaPolicyDefault_Issue_1481 deals with fixing/testing
// Panic runtime error in 3.43.0 on okta_policy_mfa_default resource #1481
// https://github.com/okta/terraform-provider-okta/issues/1481
func TestAccResourceOktaMfaPolicyDefault_Issue_1481(t *testing.T) {
	mgr := newFixtureManager(policyMfaDefault, t.Name())
	config := `
resource "okta_policy_mfa_default" "test" {
  is_oie = true
   
  okta_password = {
    enroll = "REQUIRED"
  }
  okta_email = {
    enroll = "REQUIRED"
  }
  okta_verify = {
    enroll = "OPTIONAL"
  }
}`
	resourceName := fmt.Sprintf("%s.test", policyMfaDefault)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "okta_password.enroll", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "okta_email.enroll", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "okta_verify.enroll", "OPTIONAL"),
				),
			},
			{
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "okta_password.enroll", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "okta_email.enroll", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "okta_verify.enroll", "OPTIONAL"),
				),
			},
		},
	})
}
