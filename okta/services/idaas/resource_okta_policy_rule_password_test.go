package idaas_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
	"github.com/okta/terraform-provider-okta/sdk"
)

func TestAccResourceOktaPolicyRulePassword_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSPolicyRulePassword, t.Name())
	config := testOktaPolicyRulePassword(mgr.Seed)
	updatedConfig := testOktaPolicyRulePasswordUpdated(mgr.Seed)
	resourceName := acctest.BuildResourceFQN(resources.OktaIDaaSPolicyRulePassword, mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkRuleDestroy(resources.OktaIDaaSPolicyRulePassword),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusInactive),
					resource.TestCheckResourceAttr(resourceName, "password_change", "DENY"),
					resource.TestCheckResourceAttr(resourceName, "password_reset", "DENY"),
					resource.TestCheckResourceAttr(resourceName, "password_unlock", "ALLOW"),
				),
			},
		},
	})
}

// Testing the logic that errors when an invalid priority is provided
func TestAccResourceOktaPolicyRulePassword_priorityError(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSPolicyRulePassword, t.Name())
	config := testOktaPolicyRulePriorityError(mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkRuleDestroy(resources.OktaIDaaSPolicyRulePassword),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("provided priority was not valid, got: 999, API responded with: 1. See schema for attribute details"),
			},
		},
	})
}

// Testing the successful setting of priority
func TestAccResourceOktaPolicyRulePassword_priority(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSPolicyRulePassword, t.Name())
	config := testOktaPolicyRulePriority(mgr.Seed)
	resourceName := acctest.BuildResourceFQN(resources.OktaIDaaSPolicyRulePassword, mgr.Seed)
	name := acctest.BuildResourceName(mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkRuleDestroy(resources.OktaIDaaSPolicyRulePassword),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "priority", "1"),
				),
			},
		},
	})
}

func ensureRuleExists(resourceName string) resource.TestCheckFunc {
	if os.Getenv("OKTA_VCR_TF_ACC") == "play" {
		return func(s *terraform.State) error {
			return nil
		}
	}
	return func(s *terraform.State) error {
		missingErr := fmt.Errorf("resource not found: %s", resourceName)
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return missingErr
		}

		policyID := rs.Primary.Attributes["policy_id"]
		exist, err := doesRuleExistsUpstream(policyID, rs.Primary.ID)
		if err != nil {
			return err
		} else if !exist {
			return missingErr
		}

		return nil
	}
}

func checkRuleDestroy(ruleType string) func(*terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != ruleType {
				continue
			}

			policyID := rs.Primary.Attributes["policy_id"]
			exists, err := doesRuleExistsUpstream(policyID, rs.Primary.ID)
			if err != nil {
				return err
			}

			if exists {
				return fmt.Errorf("rule still exists, ID: %s, PolicyID: %s", rs.Primary.ID, policyID)
			}
		}
		return nil
	}
}

func doesRuleExistsUpstream(policyID, ruleID string) (bool, error) {
	client := iDaaSAPIClientForTestUtil.OktaSDKSupplementClient()
	rule, resp, err := client.GetPolicyRule(context.Background(), policyID, ruleID)
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return rule.Id != "", nil
}

func testOktaPolicyRulePassword(rInt int) string {
	name := acctest.BuildResourceName(rInt)

	return fmt.Sprintf(`
data "okta_default_policy" "default-%d" {
	type = "%s"
}

resource "%s" "%s" {
	policy_id = "${data.okta_default_policy.default-%d.id}"
	name     = "%s"
	status   = "ACTIVE"
}
`, rInt, sdk.PasswordPolicyType, resources.OktaIDaaSPolicyRulePassword, name, rInt, name)
}

func testOktaPolicyRulePriority(rInt int) string {
	name := acctest.BuildResourceName(rInt)

	return fmt.Sprintf(`
data "okta_default_policy" "default-%d" {
	type = "%s"
}

resource "%s" "%s" {
	policy_id = "${data.okta_default_policy.default-%d.id}"
	name     = "%s"
	priority = 1
	status   = "ACTIVE"
}
`, rInt, sdk.PasswordPolicyType, resources.OktaIDaaSPolicyRulePassword, name, rInt, name)
}

func testOktaPolicyRulePriorityError(rInt int) string {
	name := acctest.BuildResourceName(rInt)

	return fmt.Sprintf(`
data "okta_default_policy" "default-%d" {
	type = "%s"
}

resource "%s" "%s" {
	policy_id = "${data.okta_default_policy.default-%d.id}"
	name     = "%s"
	priority = 999
	status   = "ACTIVE"
	users_included  = ["00ustguf78owmG7Rt1d7"]
	users_excluded  = ["00urzse61ohS6KPfT1d7"]
	groups_included = ["00gwxsozqariU272g1d7"]
	groups_excluded = ["00gwxstmy6w36z1dZ1d7"]
}
`, rInt, sdk.PasswordPolicyType, resources.OktaIDaaSPolicyRulePassword, name, rInt, name)
}

// TestAccResourceOktaPolicyRulePassword_sspr tests the SSPR requirement fields added in GH-2559,
// including method_constraints (otp/google_otp), primary_methods, step_up_enabled, step_up_methods,
// and switching between LEGACY and AUTH_POLICY access control.
func TestAccResourceOktaPolicyRulePassword_sspr(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSPolicyRulePassword, t.Name())
	// Step 1: LEGACY with method_constraints (config loaded from examples fixture file)
	config := mgr.GetFixtures("sspr_with_method_constraints.tf", t)
	// Step 2: LEGACY without method_constraints (config loaded from examples fixture file)
	updatedConfig := mgr.GetFixtures("sspr_no_method_constraints.tf", t)
	// Step 3: Switch to AUTH_POLICY (config loaded from examples fixture file)
	authPolicyConfig := mgr.GetFixtures("sspr_auth_policy.tf", t)
	// Step 4: LEGACY with method_constraints and step_up_methods
	stepUpMethodsConfig := mgr.GetFixtures("sspr_with_step_up_methods.tf", t)
	resourceName := acctest.BuildResourceFQN(resources.OktaIDaaSPolicyRulePassword, mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkRuleDestroy(resources.OktaIDaaSPolicyRulePassword),
		Steps: []resource.TestStep{
			{
				// LEGACY + method_constraints
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "password_reset_access_control", "LEGACY"),
					resource.TestCheckResourceAttr(resourceName, "network_connection", "ZONE"),
					resource.TestCheckResourceAttr(resourceName, "password_reset_requirement.0.step_up_enabled", "true"),
					resource.TestCheckTypeSetElemAttr(resourceName, "password_reset_requirement.0.primary_methods.*", "otp"),
					resource.TestCheckResourceAttr(resourceName, "password_reset_requirement.0.method_constraints.0.method", "otp"),
					resource.TestCheckResourceAttr(resourceName, "users_included.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "users_excluded.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "groups_included.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "groups_excluded.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_includes.#", "1"),
				),
			},
			{
				// LEGACY without method_constraints – verifies removing constraints works
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "password_reset_access_control", "LEGACY"),
					resource.TestCheckResourceAttr(resourceName, "password_reset_requirement.0.step_up_enabled", "true"),
					resource.TestCheckTypeSetElemAttr(resourceName, "password_reset_requirement.0.primary_methods.*", "push"),
					resource.TestCheckTypeSetElemAttr(resourceName, "password_reset_requirement.0.primary_methods.*", "sms"),
					resource.TestCheckResourceAttr(resourceName, "password_reset_requirement.0.method_constraints.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "users_included.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "users_excluded.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "groups_included.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "groups_excluded.#", "1"),
				),
			},
			{
				// AUTH_POLICY – SSPR delegated to authentication policy rules
				Config: authPolicyConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "password_reset_access_control", "AUTH_POLICY"),
					resource.TestCheckResourceAttr(resourceName, "users_included.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "users_excluded.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "groups_included.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "groups_excluded.#", "1"),
				),
			},
			{
				// LEGACY with method_constraints and step_up_methods=security_question
				Config: stepUpMethodsConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "password_reset_access_control", "LEGACY"),
					resource.TestCheckResourceAttr(resourceName, "password_reset_requirement.0.step_up_enabled", "true"),
					resource.TestCheckTypeSetElemAttr(resourceName, "password_reset_requirement.0.step_up_methods.*", "security_question"),
					resource.TestCheckTypeSetElemAttr(resourceName, "password_reset_requirement.0.primary_methods.*", "otp"),
					resource.TestCheckTypeSetElemAttr(resourceName, "password_reset_requirement.0.primary_methods.*", "email"),
					resource.TestCheckResourceAttr(resourceName, "password_reset_requirement.0.method_constraints.0.method", "otp"),
					resource.TestCheckResourceAttr(resourceName, "users_included.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "users_excluded.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "groups_included.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "groups_excluded.#", "1"),
				),
			},
		},
	})
}

// TestAccResourceOktaPolicyRulePassword_preservesUnmodeledProperties is the
// GH-2960 regression test. An update replaces the rule wholesale, so anything
// Okta stores that the schema does not model has to be read first and carried
// over, otherwise updating one managed attribute deletes it.
//
// actions.selfServicePasswordReset.settings.allowRecoveryEmailWithoutEnrollment
// is the property used here, but nothing about the fix is specific to it.
func TestAccResourceOktaPolicyRulePassword_preservesUnmodeledProperties(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSPolicyRulePassword, t.Name())
	config := testOktaPolicyRulePasswordUnlock(mgr.Seed, "DENY")
	updatedConfig := testOktaPolicyRulePasswordUnlock(mgr.Seed, "ALLOW")
	resourceName := acctest.BuildResourceFQN(resources.OktaIDaaSPolicyRulePassword, mgr.Seed)

	var policyID, ruleID string

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkRuleDestroy(resources.OktaIDaaSPolicyRulePassword),
		Steps: []resource.TestStep{
			{
				// Step 1: create the rule and capture the IDs the out-of-band
				// mutation needs.
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "password_unlock", "DENY"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok {
							return fmt.Errorf("resource not found: %s", resourceName)
						}
						ruleID = rs.Primary.ID
						policyID = rs.Primary.Attributes["policy_id"]
						return nil
					},
				),
			},
			{
				// Step 2: an admin sets a property the schema does not model,
				// then Terraform updates an unrelated managed attribute. The
				// property must still be there afterwards.
				PreConfig: func() {
					if policyID == "" || ruleID == "" {
						t.Fatalf("policy_id or rule id not captured from the previous step")
					}
					if err := setUnmodeledPasswordRuleProperty(policyID, ruleID, true); err != nil {
						t.Fatalf("failed to set the property out of band: %v", err)
					}
					set, err := unmodeledPasswordRulePropertyIsSet(policyID, ruleID)
					if err != nil {
						t.Fatalf("failed to read the rule back: %v", err)
					}
					if !set {
						t.Fatalf("the org did not persist allowRecoveryEmailWithoutEnrollment, nothing to preserve")
					}
				},
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "password_unlock", "ALLOW"),
					func(s *terraform.State) error {
						set, err := unmodeledPasswordRulePropertyIsSet(policyID, ruleID)
						if err != nil {
							return err
						}
						if !set {
							return fmt.Errorf("the update deleted allowRecoveryEmailWithoutEnrollment, it should have been carried over")
						}
						return nil
					},
				),
			},
		},
	})
}

// ssprSettingsProperty is the object holding the property this test preserves.
// The provider does not model it, so the SDK carries it in AdditionalProperties.
const ssprSettingsProperty = "settings"

// setUnmodeledPasswordRuleProperty mirrors the admin doing a GET, editing the
// JSON and PUTting it back, which is how the defect was found.
func setUnmodeledPasswordRuleProperty(policyID, ruleID string, allowRecoveryEmailWithoutEnrollment bool) error {
	ctx := context.Background()
	client := iDaaSAPIClientForTestUtil.OktaSDKClientV6().PolicyAPI

	inner, _, err := client.GetPolicyRule(ctx, policyID, ruleID).Execute()
	if err != nil {
		return fmt.Errorf("failed to get the rule: %w", err)
	}
	if inner == nil || inner.PasswordPolicyRule == nil {
		return fmt.Errorf("the rule was not returned as a password policy rule")
	}
	rule := inner.PasswordPolicyRule
	if rule.Actions == nil || rule.Actions.SelfServicePasswordReset == nil {
		return fmt.Errorf("the rule has no selfServicePasswordReset action to edit")
	}

	sspr := rule.Actions.SelfServicePasswordReset
	if sspr.AdditionalProperties == nil {
		sspr.AdditionalProperties = map[string]interface{}{}
	}
	sspr.AdditionalProperties[ssprSettingsProperty] = map[string]interface{}{
		"allowRecoveryEmailWithoutEnrollment": allowRecoveryEmailWithoutEnrollment,
	}

	payload := v6okta.PasswordPolicyRuleAsListPolicyRules200ResponseInner(rule)
	if _, _, err := client.ReplacePolicyRule(ctx, policyID, ruleID).PolicyRule(payload).Execute(); err != nil {
		return fmt.Errorf("failed to replace the rule: %w", err)
	}
	return nil
}

func unmodeledPasswordRulePropertyIsSet(policyID, ruleID string) (bool, error) {
	ctx := context.Background()
	client := iDaaSAPIClientForTestUtil.OktaSDKClientV6().PolicyAPI

	inner, _, err := client.GetPolicyRule(ctx, policyID, ruleID).Execute()
	if err != nil {
		return false, fmt.Errorf("failed to get the rule: %w", err)
	}
	if inner == nil || inner.PasswordPolicyRule == nil {
		return false, fmt.Errorf("the rule was not returned as a password policy rule")
	}
	actions := inner.PasswordPolicyRule.Actions
	if actions == nil || actions.SelfServicePasswordReset == nil {
		return false, nil
	}

	settings, ok := actions.SelfServicePasswordReset.AdditionalProperties[ssprSettingsProperty].(map[string]interface{})
	if !ok {
		return false, nil
	}
	allowed, _ := settings["allowRecoveryEmailWithoutEnrollment"].(bool)
	return allowed, nil
}

func testOktaPolicyRulePasswordUnlock(rInt int, passwordUnlock string) string {
	name := acctest.BuildResourceName(rInt)

	return fmt.Sprintf(`
data "okta_default_policy" "default-%d" {
	type = "%s"
}

resource "%s" "%s" {
	policy_id       = "${data.okta_default_policy.default-%d.id}"
	name            = "%s"
	status          = "ACTIVE"
	password_unlock = "%s"
}
`, rInt, sdk.PasswordPolicyType, resources.OktaIDaaSPolicyRulePassword, name, rInt, name, passwordUnlock)
}

func testOktaPolicyRulePasswordUpdated(rInt int) string {
	name := acctest.BuildResourceName(rInt)

	return fmt.Sprintf(`
data "okta_default_policy" "default-%d" {
	type = "%s"
}

resource "%s" "%s" {
	policy_id = "${data.okta_default_policy.default-%d.id}"
	name     = "%s"
	status   = "INACTIVE"
	password_change = "DENY"
	password_reset  = "DENY"
	password_unlock = "ALLOW"
}
`, rInt, sdk.PasswordPolicyType, resources.OktaIDaaSPolicyRulePassword, name, rInt, name)
}
