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
	"github.com/stretchr/testify/require"
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

// TestResourceOktaPolicyRulePassword_ssprRequirementRequiresAccessControl verifies
// that password_reset_requirement is rejected at plan time without
// password_reset_access_control. Previously the requirement was silently dropped
// from the request and the rule was saved without any recovery methods (OKTA-1252600).
func TestResourceOktaPolicyRulePassword_ssprRequirementRequiresAccessControl(t *testing.T) {
	requirement := []interface{}{
		map[string]interface{}{
			"primary_methods": []interface{}{"push", "email"},
			"step_up_enabled": true,
		},
	}
	tests := []struct {
		name    string
		config  map[string]interface{}
		wantErr bool
	}{
		{
			name: "requirement without access control",
			config: map[string]interface{}{
				"name":                       "test",
				"password_reset_requirement": requirement,
			},
			wantErr: true,
		},
		{
			name: "requirement with LEGACY access control",
			config: map[string]interface{}{
				"name":                          "test",
				"password_reset_access_control": "LEGACY",
				"password_reset_requirement":    requirement,
			},
		},
		{
			name: "access control without requirement",
			config: map[string]interface{}{
				"name":                          "test",
				"password_reset_access_control": "AUTH_POLICY",
			},
		},
		{
			name:   "neither set",
			config: map[string]interface{}{"name": "test"},
		},
	}
	r := idaas.ProviderResources()[resources.OktaIDaaSPolicyRulePassword]
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			diags := r.Validate(terraform.NewResourceConfigRaw(tc.config))
			if tc.wantErr {
				require.True(t, diags.HasError(), "expected a validation error")
				require.Contains(t, diags[0].Detail, "password_reset_access_control")
				return
			}
			require.False(t, diags.HasError(), "unexpected validation error: %+v", diags)
		})
	}
}

func TestFlattenSsprRequirement(t *testing.T) {
	legacy := "LEGACY"
	authPolicy := "AUTH_POLICY"
	required := true
	notRequired := false
	tests := []struct {
		name              string
		req               *v6okta.SsprRequirement
		wantAccessControl string
		wantRequirement   []interface{}
	}{
		{
			// Okta stored no requirement, so any prior state must be cleared.
			name:              "nil requirement clears both fields",
			req:               nil,
			wantAccessControl: "",
			wantRequirement:   []interface{}{},
		},
		{
			name: "LEGACY with methods",
			req: &v6okta.SsprRequirement{
				AccessControl: &legacy,
				Primary:       &v6okta.SsprPrimaryRequirement{Methods: []string{"email", "push"}},
				StepUp:        &v6okta.SsprStepUpRequirement{Required: &notRequired},
			},
			wantAccessControl: "LEGACY",
			wantRequirement: []interface{}{
				map[string]interface{}{
					"primary_methods":    []string{"email", "push"},
					"method_constraints": []interface{}{},
					"step_up_enabled":    false,
					"step_up_methods":    []string(nil),
				},
			},
		},
		{
			name: "LEGACY with method constraints and step up",
			req: &v6okta.SsprRequirement{
				AccessControl: &legacy,
				Primary: &v6okta.SsprPrimaryRequirement{
					Methods: []string{"otp"},
					MethodConstraints: []v6okta.AuthenticatorMethodConstraint{
						{
							Method:                v6okta.PtrString("otp"),
							AllowedAuthenticators: []v6okta.AuthenticatorIdentity{{Key: v6okta.PtrString("google_otp")}},
						},
					},
				},
				StepUp: &v6okta.SsprStepUpRequirement{Required: &required, Methods: []string{"security_question"}},
			},
			wantAccessControl: "LEGACY",
			wantRequirement: []interface{}{
				map[string]interface{}{
					"primary_methods": []string{"otp"},
					"method_constraints": []interface{}{
						map[string]interface{}{
							"method":                 "otp",
							"allowed_authenticators": []interface{}{"google_otp"},
						},
					},
					"step_up_enabled": true,
					"step_up_methods": []string{"security_question"},
				},
			},
		},
		{
			// AUTH_POLICY returns empty primary/stepUp; the block is omitted to avoid a permanent diff.
			name: "AUTH_POLICY with empty requirement",
			req: &v6okta.SsprRequirement{
				AccessControl: &authPolicy,
				Primary:       &v6okta.SsprPrimaryRequirement{Methods: []string{}},
				StepUp:        &v6okta.SsprStepUpRequirement{Required: &notRequired},
			},
			wantAccessControl: "AUTH_POLICY",
			wantRequirement:   []interface{}{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			accessControl, requirement := idaas.FlattenSsprRequirement(tc.req)
			require.Equal(t, tc.wantAccessControl, accessControl)
			require.Equal(t, tc.wantRequirement, requirement)
		})
	}
}
