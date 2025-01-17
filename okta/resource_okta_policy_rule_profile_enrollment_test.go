package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/sdk"
)

func TestAccResourceOktaPolicyRuleProfileEnrollment(t *testing.T) {
	mgr := newFixtureManager("resources", policyRuleProfileEnrollment, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", policyRuleProfileEnrollment)

	// NOTE: teardownConfig is a hack so that the okta_policy_profile_enrollment
	// okta_policy_rule_profile_enrollment resources are destoyed in step 2
	// before the inline hook and okta group are destroyed, e.g.
	// Error: failed to deactivate inline hook...
	// This pre-registration inline hook can't be deactivated because it is being used by a Profile Enrollment policy.
	teardownConfig := `
resource "okta_inline_hook" "test" {
  name    = "testAcc_replace_with_uuid"
  status  = "ACTIVE"
  type    = "com.okta.user.pre-registration"
  version = "1.0.3"

  channel = {
    type    = "HTTP"
    version = "1.0.0"
    uri     = "https://example.com/test2"
    method  = "POST"
  }
}

resource "okta_group" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "testing, testing"
}
`

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkRuleDestroy(policyRuleProfileEnrollment),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "unknown_user_action", "REGISTER"),
					resource.TestCheckResourceAttr(resourceName, "email_verification", "true"),
					resource.TestCheckResourceAttr(resourceName, "access", "ALLOW"),
					resource.TestCheckResourceAttr(resourceName, "enroll_authenticator_types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enroll_authenticator_types.0", "password"),
					resource.TestCheckResourceAttr(resourceName, "profile_attributes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "profile_attributes.0.name", "email"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "unknown_user_action", "REGISTER"),
					resource.TestCheckResourceAttr(resourceName, "email_verification", "true"),
					resource.TestCheckResourceAttr(resourceName, "access", "ALLOW"),
					resource.TestCheckResourceAttrSet(resourceName, "inline_hook_id"),
					resource.TestCheckResourceAttrSet(resourceName, "target_group_id"),
					resource.TestCheckResourceAttr(resourceName, "profile_attributes.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "profile_attributes.0.name", "email"),
					resource.TestCheckResourceAttr(resourceName, "profile_attributes.1.name", "mobilePhone"),
				),
			},
			{
				Config: mgr.ConfigReplace(teardownConfig),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_group.test", "name"),
				),
			},
		},
	})
}

// TestAccResourceOktaPolicyRuleProfileEnrollment_Issue1213
// re: uiSchemaId / ui_schema_id
// https://developer.okta.com/docs/reference/api/policy/#profile-enrollment-action-object
// https://github.com/okta/terraform-provider-okta/issues/1213
func TestAccResourceOktaPolicyRuleProfileEnrollment_Issue1213(t *testing.T) {
	mgr := newFixtureManager("resources", policyRuleProfileEnrollment, t.Name())
	resourceName := fmt.Sprintf("%s.test", policyRuleProfileEnrollment)
	config := `
resource "okta_policy_profile_enrollment" "test" {
  name   = "testAcc_replace_with_uuid"
  status = "ACTIVE"
}
  
resource "okta_group" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "terraform group"
}
  
resource "okta_policy_rule_profile_enrollment" "test" {
  policy_id           = okta_policy_profile_enrollment.test.id
  target_group_id     = okta_group.test.id
  unknown_user_action = "REGISTER"
  email_verification  = false
  access              = "ALLOW"
  ui_schema_id        = "uis44fio9ifOCwJAO1d7"
  profile_attributes {
    name     = "email"
    label    = "Primary Email"
    required = true
  }
  profile_attributes {
    name     = "firstName"
    label    = "First Name"
    required = true
  }
  profile_attributes {
    name     = "lastName"
    label    = "Last Name"
    required = true
  }
}`
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(appSecurePasswordStore, createDoesAppExist(sdk.NewSecurePasswordStoreApplication())),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "unknown_user_action", "REGISTER"),
					resource.TestCheckResourceAttr(resourceName, "email_verification", "false"),
					resource.TestCheckResourceAttr(resourceName, "access", "ALLOW"),
					resource.TestCheckResourceAttr(resourceName, "profile_attributes.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "profile_attributes.0.name", "email"),
					resource.TestCheckResourceAttr(resourceName, "profile_attributes.1.name", "firstName"),
					resource.TestCheckResourceAttr(resourceName, "profile_attributes.2.name", "lastName"),
					resource.TestCheckResourceAttr(resourceName, "ui_schema_id", "uis44fio9ifOCwJAO1d7"),
				),
			},
		},
	})
}
