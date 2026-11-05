package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/stretchr/testify/require"
)

func TestAccResourceOktaPolicyRuleProfileEnrollment_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSPolicyRuleProfileEnrollment, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSPolicyRuleProfileEnrollment)

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

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkRuleDestroy(resources.OktaIDaaSPolicyRuleProfileEnrollment),
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

// TestAccResourceOktaPolicyRuleProfileEnrollment_inlineHookScopes verifies the
// registration inline hook "scopes" (which registration flows run the hook) can
// be managed and are not dropped by the full-PUT update.
// https://github.com/okta/terraform-provider-okta/issues/1624
func TestAccResourceOktaPolicyRuleProfileEnrollment_inlineHookScopes(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSPolicyRuleProfileEnrollment, t.Name())
	config := mgr.GetFixtures("inline_hook_scopes.tf", t)
	updatedConfig := mgr.GetFixtures("inline_hook_scopes_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSPolicyRuleProfileEnrollment)

	// See TestAccResourceOktaPolicyRuleProfileEnrollment_crud for why the
	// teardown step is needed.
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

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkRuleDestroy(resources.OktaIDaaSPolicyRuleProfileEnrollment),
		Steps: []resource.TestStep{
			{
				// inline_hook_scopes not configured: Okta defaults to
				// SELF_SERVICE_REGISTRATION and the value is read back
				// into state without producing a diff.
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "inline_hook_id"),
					resource.TestCheckResourceAttr(resourceName, "inline_hook_scopes.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "inline_hook_scopes.*", "SELF_SERVICE_REGISTRATION"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "progressive_profiling_action", "ENABLED"),
					resource.TestCheckResourceAttr(resourceName, "inline_hook_scopes.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "inline_hook_scopes.*", "SELF_SERVICE_REGISTRATION"),
					resource.TestCheckTypeSetElemAttr(resourceName, "inline_hook_scopes.*", "PROGRESSIVE_PROFILING"),
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

func TestBuildProfileEnrollmentInlineHooks(t *testing.T) {
	both := []string{"SELF_SERVICE_REGISTRATION", "PROGRESSIVE_PROFILING"}
	existing := []*sdk.PreRegistrationInlineHook{{InlineHookId: "hook1", Scopes: both}}

	tests := []struct {
		name     string
		existing []*sdk.PreRegistrationInlineHook
		hookID   string
		scopes   []string
		expected []*sdk.PreRegistrationInlineHook
	}{
		{
			name:     "no hook configured and none on rule",
			expected: nil,
		},
		{
			name:     "no hook configured keeps hook already on rule",
			existing: existing,
			expected: existing,
		},
		{
			name:     "hook without scopes on new rule leaves scopes to the API default",
			hookID:   "hook1",
			expected: []*sdk.PreRegistrationInlineHook{{InlineHookId: "hook1"}},
		},
		{
			name:     "hook without scopes preserves scopes already on rule",
			existing: existing,
			hookID:   "hook1",
			expected: []*sdk.PreRegistrationInlineHook{{InlineHookId: "hook1", Scopes: both}},
		},
		{
			name:     "different hook does not inherit scopes",
			existing: existing,
			hookID:   "hook2",
			expected: []*sdk.PreRegistrationInlineHook{{InlineHookId: "hook2"}},
		},
		{
			name:     "configured scopes override and are sorted",
			existing: []*sdk.PreRegistrationInlineHook{{InlineHookId: "hook1", Scopes: []string{"SELF_SERVICE_REGISTRATION"}}},
			hookID:   "hook1",
			scopes:   both,
			expected: []*sdk.PreRegistrationInlineHook{{InlineHookId: "hook1", Scopes: []string{"PROGRESSIVE_PROFILING", "SELF_SERVICE_REGISTRATION"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, idaas.BuildProfileEnrollmentInlineHooks(tt.existing, tt.hookID, tt.scopes))
		})
	}
}

// TestAccResourceOktaPolicyRuleProfileEnrollment_Issue1213
// re: uiSchemaId / ui_schema_id
// https://developer.okta.com/docs/reference/api/policy/#profile-enrollment-action-object
// https://github.com/okta/terraform-provider-okta/issues/1213
func TestAccResourceOktaPolicyRuleProfileEnrollment_Issue1213(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSPolicyRuleProfileEnrollment, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSPolicyRuleProfileEnrollment)
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
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppSecurePasswordStore, createDoesAppExist(sdk.NewSecurePasswordStoreApplication())),
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
