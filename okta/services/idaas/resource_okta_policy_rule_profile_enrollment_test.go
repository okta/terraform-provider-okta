package idaas_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaPolicyRuleProfileEnrollment_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSPolicyRuleProfileEnrollment, t.Name())
	infraConfig := mgr.GetFixtures("infra.tf", t)
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	updatedScopesConfig := mgr.GetFixtures("basic_updated_scopes.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSPolicyRuleProfileEnrollment)

	// NOTE: A profile enrollment policy always has exactly one system
	// "Catch-all Rule" and the API rejects creating any other rule in it
	// ("You cannot create a default rule in a policy." / "Adding an additional
	// rule for this policy is not supported at this time."). So the test
	// creates the policy first, imports its catch-all rule into state, and then
	// exercises updates against it.
	//
	// NOTE: teardownConfig drops the rule from state with a removed block
	// (system rules can't be deleted on their own; Okta removes it with the
	// policy) and destroys the policy before the inline hook and okta group
	// are destroyed, e.g.
	// Error: failed to deactivate inline hook...
	// This pre-registration inline hook can't be deactivated because it is being used by a Profile Enrollment policy.
	teardownConfig := mgr.GetFixtures("teardown.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkRuleDestroy(resources.OktaIDaaSPolicyRuleProfileEnrollment),
		Steps: []resource.TestStep{
			{
				Config: infraConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_policy_profile_enrollment.test", "id"),
				),
			},
			{
				Config:             config,
				ResourceName:       resourceName,
				ImportState:        true,
				ImportStatePersist: true,
				ImportStateIdFunc:  profileEnrollmentCatchAllRuleImportID(t.Name(), "okta_policy_profile_enrollment.test"),
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					// With ImportStatePersist every resource in state is passed in,
					// so pick out the imported rule.
					for _, is := range s {
						if is.Ephemeral.Type != resources.OktaIDaaSPolicyRuleProfileEnrollment {
							continue
						}
						if is.Attributes["system"] != "true" {
							return fmt.Errorf("expected imported rule to be the system catch-all rule, got system=%q", is.Attributes["system"])
						}
						return nil
					}
					return fmt.Errorf("%s was not imported", resourceName)
				},
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "PROFILE_ENROLLMENT"),
					resource.TestCheckResourceAttr(resourceName, "name", "Catch-all Rule"),
					resource.TestCheckResourceAttr(resourceName, "system", "true"),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.access", "ALLOW"),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.unknown_user_action", "REGISTER"),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.progressive_profiling_action", "ENABLED"),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.allowed_identifiers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.allowed_identifiers.0", "login"),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.target_group_ids.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "actions.profile_enrollment.target_group_ids.0", "okta_group.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.activation_requirements.email_verification", "true"),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.pre_registration_inline_hooks.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "actions.profile_enrollment.pre_registration_inline_hooks.0.inline_hook_id", "okta_inline_hook.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.pre_registration_inline_hooks.0.scopes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.pre_registration_inline_hooks.0.scopes.0", "SELF_SERVICE_REGISTRATION"),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.profile_attributes.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.profile_attributes.0.name", "email"),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.profile_attributes.1.name", "firstName"),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.profile_attributes.2.name", "lastName"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.activation_requirements.email_verification", "false"),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.profile_attributes.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.profile_attributes.3.name", "mobilePhone"),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.profile_attributes.3.required", "false"),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.pre_registration_inline_hooks.0.scopes.#", "2"),
				),
			},
			{
				// Narrowing scopes must remove SELF_SERVICE_REGISTRATION upstream, not keep the prior value.
				Config: updatedScopesConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.pre_registration_inline_hooks.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "actions.profile_enrollment.pre_registration_inline_hooks.0.inline_hook_id", "okta_inline_hook.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.pre_registration_inline_hooks.0.scopes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.pre_registration_inline_hooks.0.scopes.0", "PROGRESSIVE_PROFILING"),
				),
			},
			{
				Config: teardownConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_group.test", "name"),
				),
			},
		},
	})
}

// profileEnrollmentCatchAllRuleImportID returns an ImportStateIdFunc that
// looks up the system catch-all rule of the profile enrollment policy at
// policyResourceName and returns its "{policy_id}/{rule_id}" import ID. When
// VCR is enabled the lookup goes through the provider's VCR client so it is
// recorded/played with the rest of the test.
func profileEnrollmentCatchAllRuleImportID(testName, policyResourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[policyResourceName]
		if !ok {
			return "", fmt.Errorf("failed to find %s in state", policyResourceName)
		}
		policyID := rs.Primary.ID

		client := acctest.VCRIDaaSClientForTest(testName)
		if client == nil {
			client = iDaaSAPIClientForTestUtil
		}
		rules, _, err := client.OktaSDKClientV8().PolicyAPI.ListPolicyRules(context.Background(), policyID).Execute()
		if err != nil {
			return "", fmt.Errorf("failed to list rules for policy %s: %w", policyID, err)
		}
		for _, rule := range rules {
			if r := rule.ProfileEnrollmentPolicyRule; r != nil && r.GetSystem() {
				return fmt.Sprintf("%s/%s", policyID, r.GetId()), nil
			}
		}
		return "", fmt.Errorf("policy %s has no system catch-all profile enrollment rule", policyID)
	}
}

// TestAccResourceOktaPolicyRuleProfileEnrollment_Issue1213
// re: uiSchemaId / ui_schema_id
// https://developer.okta.com/docs/reference/api/policy/#profile-enrollment-action-object
// https://github.com/okta/terraform-provider-okta/issues/1213
func TestAccResourceOktaPolicyRuleProfileEnrollment_Issue1213(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSPolicyRuleProfileEnrollment, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSPolicyRuleProfileEnrollment)
	infraConfig := mgr.GetFixtures("ui_schema_infra.tf", t)
	config := mgr.GetFixtures("ui_schema_id.tf", t)
	// See TestAccResourceOktaPolicyRuleProfileEnrollment_crud: the catch-all
	// rule can't be created or deleted, so it is imported and later dropped
	// from state. The UI schema is created by the test rather than hard coded
	// so the test isn't tied to one org.
	teardownConfig := mgr.GetFixtures("ui_schema_teardown.tf", t)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkRuleDestroy(resources.OktaIDaaSPolicyRuleProfileEnrollment),
		Steps: []resource.TestStep{
			{
				Config: infraConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_policy_profile_enrollment.test", "id"),
					resource.TestCheckResourceAttrSet("okta_ui_schema.test", "id"),
				),
			},
			{
				Config:             config,
				ResourceName:       resourceName,
				ImportState:        true,
				ImportStatePersist: true,
				ImportStateIdFunc:  profileEnrollmentCatchAllRuleImportID(t.Name(), "okta_policy_profile_enrollment.test"),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.unknown_user_action", "REGISTER"),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.activation_requirements.email_verification", "false"),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.access", "ALLOW"),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.profile_attributes.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.profile_attributes.0.name", "email"),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.profile_attributes.1.name", "firstName"),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.profile_attributes.2.name", "lastName"),
					resource.TestCheckResourceAttr(resourceName, "actions.profile_enrollment.progressive_profiling_action", "DISABLED"),
					resource.TestCheckResourceAttrPair(resourceName, "actions.profile_enrollment.ui_schema_id", "okta_ui_schema.test", "id"),
				),
			},
			{
				Config: teardownConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_ui_schema.test", "id"),
				),
			},
		},
	})
}
