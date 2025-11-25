package idaas_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
)

// TestAccResourceOktaAppSignOnPolicyRule_crud can flap when all the tests are
// run in the harness but rarely fails running individually.
func TestAccResourceOktaAppSignOnPolicyRule_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppSignOnPolicyRule)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSignOnPolicyRule, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkAppSignOnPolicyRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttrSet(resourceName, "priority"),
					resource.TestCheckResourceAttr(resourceName, "access", "ALLOW"),
					resource.TestCheckResourceAttr(resourceName, "factor_mode", "2FA"),
					resource.TestCheckResourceAttr(resourceName, "groups_excluded.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "groups_included.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "device_assurances_included.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "user_types_excluded.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "user_types_included.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "users_excluded.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "users_included.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "network_includes.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "network_excludes.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "network_connection", "ANYWHERE"),
					resource.TestCheckResourceAttr(resourceName, "constraints.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "re_authentication_frequency", "PT2H"),
					resource.TestCheckResourceAttr(resourceName, "inactivity_period", "PT1H"),
					resource.TestCheckResourceAttr(resourceName, "risk_score", "LOW"),
					resource.TestCheckResourceAttr(resourceName, "platform_include.#", "1"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)+"_updated"),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttrSet(resourceName, "priority"),
					resource.TestCheckResourceAttr(resourceName, "access", "ALLOW"),
					resource.TestCheckResourceAttr(resourceName, "factor_mode", "2FA"),
					resource.TestCheckResourceAttr(resourceName, "groups_excluded.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "groups_included.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "device_assurances_included.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "user_types_excluded.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "user_types_included.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "users_excluded.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "users_included.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "network_includes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_excludes.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "network_connection", "ZONE"),
					resource.TestCheckResourceAttr(resourceName, "platform_include.#", "5"),
					resource.TestCheckResourceAttr(resourceName, "re_authentication_frequency", "PT43800H"),
					resource.TestCheckResourceAttr(resourceName, "inactivity_period", "PT2H"),
					resource.TestCheckResourceAttr(resourceName, "type", "ASSURANCE"),
					resource.TestCheckResourceAttr(resourceName, "constraints.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "risk_score", "MEDIUM"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("failed to find app sign on policy rule %s", resourceName)
					}
					return fmt.Sprintf("%s/%s", rs.Primary.Attributes["policy_id"], rs.Primary.Attributes["id"]), nil
				},
				ImportStateCheck: func(s []*terraform.InstanceState) (err error) {
					if len(s) != 1 {
						err = errors.New("failed to import into resource into state")
						return
					}

					id := s[0].Attributes["id"]
					if strings.Contains(id, "@") {
						err = fmt.Errorf("resource id incorrectly set, %s", id)
					}
					return
				},
			},
		},
	})
}

// TestAccResourceOktaAppSignOnPolicyRule_Issue_1242_possession_constraint
// https://github.com/okta/terraform-provider-okta/issues/1242
// Operator had a typo in the constraint, possession and not possession. We'll
// still keep this ACC.
func TestAccResourceOktaAppSignOnPolicyRule_Issue_1242_possession_constraint(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSignOnPolicyRule, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppSignOnPolicyRule)
	constraints := []interface{}{
		map[string]interface{}{
			"knowledge": map[string]interface{}{
				"reauthenticateIn": "PT43800H",
				"required":         false,
				"types":            []string{"password"},
			},
			"possession": map[string]interface{}{
				"required":    false,
				"deviceBound": "REQUIRED",
			},
		},
	}
	config := `
resource "okta_app_signon_policy" "test" {
	name        = "testAcc_replace_with_uuid"
	description = "Test App Signon Policy with updated Okta TF Provider"
}
resource "okta_app_signon_policy_rule" "test" {
	policy_id                   = okta_app_signon_policy.test.id
	name                        = "Require MFA_replace_with_uuid"
	access                      = "ALLOW"
	re_authentication_frequency = "PT43800H"
	constraints = [
		jsonencode({
			knowledge = {
				reauthenticateIn = "PT43800H"
				types            = ["password"],
				required         = false
			}
			possession = {
			  deviceBound = "REQUIRED"
			  required    = false
			}
	  })
	]
}`

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		CheckDestroy:             nil,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceNameWithPrefix("Require MFA", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "access", "ALLOW"),
					resource.TestCheckResourceAttr(resourceName, "re_authentication_frequency", "PT43800H"),
					// Note: validateOktaAppSignonPolicyRuleConstraintsAreSet no longer works correctly with the addition of required
					validateOktaAppSignonPolicyRuleConstraintsAreSet(resourceName, constraints),
				),
			},
		},
	})
}

// TestAccResourceOktaAppSignOnPolicyRule_Issue_1245_existing_default_rule
// https://github.com/okta/terraform-provider-okta/issues/1245
// This ACC was used to find and fix the issues with importing then interacting
// with the default rule of an okta_app_signon_policy
func TestAccResourceOktaAppSignOnPolicyRule_Issue_1245_existing_default_rule(t *testing.T) {
	resourceName := "okta_app_signon_policy_rule.test"
	policyName := "okta_app_signon_policy.test"
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSignOnPolicyRule, t.Name())
	baseConfig := `
resource "okta_app_signon_policy" "test" {
	name        = "testAcc_replace_with_uuid"
	description = "Test App Signon Policy with updated Okta TF Provider"
}`

	importConfig := `
resource "okta_app_signon_policy_rule" "test" {
	name                        = "Catch-all Rule"
	policy_id                   = okta_app_signon_policy.test.id
	constraints = [
		jsonencode({
			possession = {
			  deviceBound = "REQUIRED"
			  required    = false
			}
	  })
	]
}`
	step4Config := `
resource "okta_app_signon_policy_rule" "test" {
	name                        = "Catch-all Rule"
	policy_id                   = okta_app_signon_policy.test.id
	inactivity_period           = "PT1H"
	re_authentication_frequency = "PT2H"
	constraints = [
		jsonencode({
			possession = {
			  deviceBound = "REQUIRED"
			  required    = false
			}
	  })
	]
}`
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		CheckDestroy:             nil,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("okta_app_signon_policy.test", "name", acctest.BuildResourceName(mgr.Seed)),
				),
			},
			{
				ResourceName:       resourceName,
				ImportState:        true,
				ImportStatePersist: true,
				Config:             mgr.ConfigReplace(fmt.Sprintf("%s\n%s", baseConfig, importConfig)),
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					policy, ok := s.RootModule().Resources[policyName]
					if !ok {
						return "", fmt.Errorf("failed to find app sign on policy%s", policyName)
					}

					policyID := policy.Primary.Attributes["id"]
					client := iDaaSAPIClientForTestUtil.OktaSDKClientV2()

					rules, _, err := client.Policy.ListPolicyRules(context.Background(), policyID)
					if err != nil {
						return "", err
					}
					if len(rules) != 1 {
						return "", fmt.Errorf("at this point, policy %q should only have one rule, its default rule", policyID)
					}
					return fmt.Sprintf("%s/%s", policyID, rules[0].Id), nil
				},
			},
			{
				Config: mgr.ConfigReplace(fmt.Sprintf("%s\n%s", baseConfig, importConfig)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Catch-all Rule"),
					resource.TestCheckResourceAttr(resourceName, "system", "true"),
				),
			},
			{
				// make sure we can update some non-conditions arguments like inactivity period and reauth freq
				Config: mgr.ConfigReplace(fmt.Sprintf("%s\n%s", baseConfig, step4Config)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Catch-all Rule"),
					resource.TestCheckResourceAttr(resourceName, "system", "true"),
					resource.TestCheckResourceAttr(resourceName, "inactivity_period", "PT1H"),
					resource.TestCheckResourceAttr(resourceName, "re_authentication_frequency", "PT2H"),
				),
			},
		},
	})
}

func checkAppSignOnPolicyRuleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resources.OktaIDaaSAppSignOnPolicyRule {
			continue
		}
		client := iDaaSAPIClientForTestUtil.OktaSDKSupplementClient()
		rule, resp, err := client.GetAppSignOnPolicyRule(context.Background(), rs.Primary.Attributes["policy_id"], rs.Primary.ID)
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		} else if err != nil {
			return err
		}
		if rule != nil {
			return fmt.Errorf("app sign-on policy rule still exists, ID: %s, PolicyID: %s", rs.Primary.ID, rs.Primary.Attributes["policy_id"])
		}
		return nil
	}
	return nil
}

func validateOktaAppSignonPolicyRuleConstraintsAreSet(rule string, expectedConstraints []interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		missingErr := fmt.Errorf("rule not found: %s", rule)
		ruleRS, ok := s.RootModule().Resources[rule]
		if !ok {
			return missingErr
		}
		ruleID := ruleRS.Primary.Attributes["id"]
		policyID := ruleRS.Primary.Attributes["policy_id"]

		client := iDaaSAPIClientForTestUtil.OktaSDKSupplementClient()
		r, _, err := client.GetAppSignOnPolicyRule(context.Background(), policyID, ruleID)
		if err != nil {
			return fmt.Errorf("API: to get policy/rule %q/%q, err: %+v", policyID, ruleID, err)
		}
		constraintsJSON, err := json.Marshal(r.Actions.AppSignOn.VerificationMethod.Constraints)
		if err != nil {
			return fmt.Errorf("unable to marshal constraints, err: %+v", err)
		}
		var gotConstraints []interface{}
		err = json.Unmarshal([]byte(constraintsJSON), &gotConstraints)
		if err != nil {
			return fmt.Errorf("unable to unmarshal constraints, err: %+v", err)
		}
		if reflect.DeepEqual(expectedConstraints, gotConstraints) {
			// object equivelence, ok
			return nil
		}
		expectedJSON, _ := json.Marshal(expectedConstraints)

		// reserialize to absolutely generatic objects
		var _expected []interface{}
		var _got []interface{}
		_ = json.Unmarshal([]byte(expectedJSON), &_expected)
		_ = json.Unmarshal([]byte(constraintsJSON), &_got)
		if reflect.DeepEqual(_expected, _got) {
			// absolute object equivelence, ok
			return nil
		}

		// last attempt of string equivelence of JSON is brittle comparing the
		// string literal of the two constraints
		if !reflect.DeepEqual(expectedJSON, constraintsJSON) {
			return fmt.Errorf("expected constraints:\n%s\ngot:\n%s", expectedJSON, constraintsJSON)
		}

		return nil
	}
}

func TestAccResourceOktaAppSignOnPolicyRule_default_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppSignOnPolicyRule)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSignOnPolicyRule, t.Name())
	config := mgr.GetFixtures("default.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkAppSignOnPolicyRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttrSet(resourceName, "priority"),
					resource.TestCheckResourceAttr(resourceName, "access", "ALLOW"),
					resource.TestCheckResourceAttr(resourceName, "factor_mode", "2FA"),
					resource.TestCheckResourceAttr(resourceName, "groups_excluded.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "groups_included.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "device_assurances_included.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "user_types_excluded.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "user_types_included.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "users_excluded.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "users_included.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "network_includes.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "network_excludes.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "network_connection", "ANYWHERE"),
					resource.TestCheckResourceAttr(resourceName, "constraints.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "re_authentication_frequency", "PT2H"),
					resource.TestCheckResourceAttr(resourceName, "inactivity_period", "PT1H"),
					resource.TestCheckResourceAttr(resourceName, "risk_score", "ANY"),
				),
			},
		},
	})
}

func TestAccResourceOktaAppSignOnPolicyRule_os_expression_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppSignOnPolicyRule)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSignOnPolicyRule, t.Name())
	config := mgr.GetFixtures("os_expression.tf", t)
	updatedConfig := mgr.GetFixtures("os_expression_updated.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkAppSignOnPolicyRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test1"),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "platform_include.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "platform_include.0.os_expression", ""),
					resource.TestCheckResourceAttr(resourceName, "platform_include.0.os_type", "OTHER"),
					resource.TestCheckResourceAttr(resourceName, "platform_include.0.type", "DESKTOP"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test1"),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "platform_include.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "platform_include.0.os_expression", ""),
					resource.TestCheckResourceAttr(resourceName, "platform_include.0.os_type", "IOS"),
					resource.TestCheckResourceAttr(resourceName, "platform_include.0.type", "MOBILE"),
				),
			},
		},
	})
}

func TestAccResourceOktaAppSignOnPolicyRule_AUTH_METHOD_CHAIN(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppSignOnPolicyRule)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSignOnPolicyRule, t.Name())
	config := mgr.GetFixtures("auth_chain_method.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkAppSignOnPolicyRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test2"),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "chains.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "chains.0", "{\"authenticationMethods\":[{\"key\":\"okta_password\",\"method\":\"password\"}],\"next\":[{\"authenticationMethods\":[{\"key\":\"okta_email\",\"method\":\"email\"}]}]}"),
				),
			},
		},
	})
}

func TestAccResourceOktaAppSignOnPolicyRule_ReauthenticationFrequency(t *testing.T) {
	resourceName1 := fmt.Sprintf("%s.test_with_reauthenticate_in_chains_only", resources.OktaIDaaSAppSignOnPolicyRule)
	resourceName2 := fmt.Sprintf("%s.test_with_re_authentication_frequency_only", resources.OktaIDaaSAppSignOnPolicyRule)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSignOnPolicyRule, t.Name())
	config := mgr.GetFixtures("reauthentication.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkAppSignOnPolicyRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName1, "chains.0", checkReauthenticateInChains),
					resource.TestCheckResourceAttrWith(resourceName1, "chains.1", checkReauthenticateInChains),
					resource.TestCheckResourceAttr(resourceName2, "re_authentication_frequency", "PT2H10M"),
				),
			},
		},
	})
}

// TestAccResourceOktaAppSignOnPolicyRule_priority_concurrency is a test to
// verify that the policy-specific mutex locking prevents concurrent modification
// issues when creating/updating multiple rules with priorities. This test ensures
// that the Okta API's automatic priority shifting behavior works correctly with
// the provider's mutex implementation.
func TestAccResourceOktaAppSignOnPolicyRule_priority_concurrency(t *testing.T) {
	numRules := 10
	testPolicyRules := make([]string, numRules)
	// Test setup makes each policy rule dependent on the one before it.
	for i := 0; i < numRules; i++ {
		dependsOn := i - 1
		testPolicyRules[i] = testAppSignOnPolicyRule(i, dependsOn)
	}
	config := fmt.Sprintf(`
resource "okta_app_signon_policy" "test" {
	name        = "testAcc_replace_with_uuid"
	description = "Test App Signon Policy for Priority Concurrency"
}
%s`, strings.Join(testPolicyRules, ""))

	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSignOnPolicyRule, t.Name())
	acctest.OktaResourceTest(
		t, resource.TestCase{
			PreCheck:                 acctest.AccPreCheck(t),
			ErrorCheck:               testAccErrorChecks(t),
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
			CheckDestroy:             checkAppSignOnPolicyRuleDestroy,
			Steps: []resource.TestStep{
				{
					Config: mgr.ConfigReplace(config),
					Check: resource.ComposeTestCheckFunc(
						// Check if all rules were created successfully without 500 errors
						resource.TestCheckResourceAttr("okta_app_signon_policy_rule.test_00", "name", "Test App Sign On Policy Rule 00"),
						resource.TestCheckResourceAttr("okta_app_signon_policy_rule.test_00", "priority", "1"),
						resource.TestCheckResourceAttr("okta_app_signon_policy_rule.test_09", "name", "Test App Sign On Policy Rule 09"),
						resource.TestCheckResourceAttr("okta_app_signon_policy_rule.test_09", "priority", "10"),
					),
				},
			},
		})
}

// Helper function to generate app sign-on policy rule configurations
func testAppSignOnPolicyRule(num, dependsOn int) string {
	var dependsOnStr string
	if dependsOn >= 0 {
		dependsOnStr = fmt.Sprintf("depends_on = [okta_app_signon_policy_rule.test_%02d]", dependsOn)
	}
	return fmt.Sprintf(`
resource "okta_app_signon_policy_rule" "test_%02d" {
	policy_id = okta_app_signon_policy.test.id
	name      = "Test App Sign On Policy Rule %02d"
	priority  = %d
	access    = "ALLOW"
	%s
}`,
		num, num, num+1, dependsOnStr)
}

func checkReauthenticateInChains(value string) error {
	if strings.Contains(value, `"reauthenticateIn":"PT43800H"`) {
		return nil
	}
	return fmt.Errorf("chains does not contain expected value")
}
