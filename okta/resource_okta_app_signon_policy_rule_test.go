package okta

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
)

// TODO unable to run the test due to conflict providerFactories between plugin and framework
func TestAccResourceOktaAppSignOnPolicyRule(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", appSignOnPolicyRule)
	mgr := newFixtureManager("resources", appSignOnPolicyRule, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:                 testAccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: testAccMergeProvidersFactories,
		CheckDestroy:             checkAppSignOnPolicyRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
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
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)+"_updated"),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
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

// TestAccResourceOktaAppSignOnPolicyRule_Issue_1275_default_conditions
// https://github.com/okta/terraform-provider-okta/issues/1275
func TestAccResourceOktaAppSignOnPolicyRule_Issue_1275_default_conditions(t *testing.T) {
	mgr := newFixtureManager("resources", appSignOnPolicyRule, t.Name())
	resourceName := fmt.Sprintf("%s.test", appSignOnPolicyRule)
	config := `
resource "okta_app_signon_policy" "test" {
	name        = "testAcc_replace_with_uuid"
	description = "Test App Signon Policy with updated Okta TF Provider"
}
resource "okta_app_signon_policy_rule" "test" {
	name                        = "Catch-all Rule_replace_with_uuid"
	policy_id                   = okta_app_signon_policy.test.id
	access                      = "ALLOW"
	priority                    = 89
	inactivity_period           = "PT0S"
	re_authentication_frequency = "PT12H"
	constraints = [
	  jsonencode({
		"possession" : {
		  "deviceBound" : "REQUIRED"
		}
	  })
	]
}`
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceNameWithPrefix("Catch-all Rule", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "priority", "89"),
					resource.TestCheckResourceAttr(resourceName, "inactivity_period", "PT0S"),
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

// TestAccResourceOktaAppSignOnPolicyRule_Issue_1242_possesion_constraint
// https://github.com/okta/terraform-provider-okta/issues/1242
// Operator had a typo in the constraint, posession and not possession. We'll
// still keep this ACC.
func TestAccResourceOktaAppSignOnPolicyRule_Issue_1242_possesion_constraint(t *testing.T) {
	mgr := newFixtureManager("resources", appSignOnPolicyRule, t.Name())
	resourceName := fmt.Sprintf("%s.test", appSignOnPolicyRule)
	constraints := []interface{}{
		map[string]interface{}{
			"knowledge": map[string]interface{}{
				"reauthenticateIn": "PT43800H",
				"types":            []string{"password"},
			},
			"possession": map[string]interface{}{
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
				types            = ["password"]
			}
			possession = {
			  deviceBound = "REQUIRED"
			}
	  })
	]
}`
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceNameWithPrefix("Require MFA", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "access", "ALLOW"),
					resource.TestCheckResourceAttr(resourceName, "re_authentication_frequency", "PT43800H"),
					validateOktaAppSignonPolicyRuleConstraintsAreSet(resourceName, constraints),
				),
			},
		},
	})
}

func checkAppSignOnPolicyRuleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != appSignOnPolicyRule {
			continue
		}
		client := sdkSupplementClientForTest()
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

func validateOktaAppSignonPolicyRuleConstraintsAreSet(rule string, constraints []interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		missingErr := fmt.Errorf("rule not found: %s", rule)
		ruleRS, ok := s.RootModule().Resources[rule]
		if !ok {
			return missingErr
		}
		ruleID := ruleRS.Primary.Attributes["id"]
		policyID := ruleRS.Primary.Attributes["policy_id"]

		client := sdkSupplementClientForTest()
		r, _, err := client.GetAppSignOnPolicyRule(context.Background(), policyID, ruleID)
		if err != nil {
			return fmt.Errorf("API: to get policy/rule %q/%q, err: %+v", policyID, ruleID, err)
		}
		constraintsJSON, err := json.Marshal(r.Actions.AppSignOn.VerificationMethod.Constraints)
		if err != nil {
			return fmt.Errorf("unable to marshal constraints, err: %+v", err)
		}
		var _constraints []interface{}
		err = json.Unmarshal([]byte(constraintsJSON), &_constraints)
		if err != nil {
			return fmt.Errorf("unable to unmarshal constraints, err: %+v", err)
		}
		expectedJSON, _ := json.Marshal(constraints)
		// NOTE: this could be brittle comparing the string literal of the two constraints
		if !reflect.DeepEqual(expectedJSON, constraintsJSON) {
			return fmt.Errorf("expected constraints:\n%s\ngot:\n%s", expectedJSON, constraintsJSON)
		}

		return nil
	}
}
