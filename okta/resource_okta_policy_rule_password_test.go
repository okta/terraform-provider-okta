package okta

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/sdk"
)

func TestAccResourceOktaPolicyRulePassword_crud(t *testing.T) {
	mgr := newFixtureManager("", policyRulePassword, t.Name())
	config := testOktaPolicyRulePassword(mgr.Seed)
	updatedConfig := testOktaPolicyRulePasswordUpdated(mgr.Seed)
	resourceName := buildResourceFQN(policyRulePassword, mgr.Seed)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkRuleDestroy(policyRulePassword),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusInactive),
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
	mgr := newFixtureManager("", policyRulePassword, t.Name())
	config := testOktaPolicyRulePriorityError(mgr.Seed)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkRuleDestroy(policyRulePassword),
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
	mgr := newFixtureManager("", policyRulePassword, t.Name())
	config := testOktaPolicyRulePriority(mgr.Seed)
	resourceName := buildResourceFQN(policyRulePassword, mgr.Seed)
	name := buildResourceName(mgr.Seed)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkRuleDestroy(policyRulePassword),
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
	client := sdkSupplementClientForTest()
	rule, resp, err := client.GetPolicyRule(context.Background(), policyID, ruleID)
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return rule.Id != "", nil
}

func testOktaPolicyRulePassword(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
data "okta_default_policy" "default-%d" {
	type = "%s"
}

resource "%s" "%s" {
	policy_id = "${data.okta_default_policy.default-%d.id}"
	name     = "%s"
	status   = "ACTIVE"
}
`, rInt, sdk.PasswordPolicyType, policyRulePassword, name, rInt, name)
}

func testOktaPolicyRulePriority(rInt int) string {
	name := buildResourceName(rInt)

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
`, rInt, sdk.PasswordPolicyType, policyRulePassword, name, rInt, name)
}

func testOktaPolicyRulePriorityError(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
data "okta_default_policy" "default-%d" {
	type = "%s"
}

resource "%s" "%s" {
	policy_id = "${data.okta_default_policy.default-%d.id}"
	name     = "%s"
	priority = 999
	status   = "ACTIVE"
}
`, rInt, sdk.PasswordPolicyType, policyRulePassword, name, rInt, name)
}

func testOktaPolicyRulePasswordUpdated(rInt int) string {
	name := buildResourceName(rInt)

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
`, rInt, sdk.PasswordPolicyType, policyRulePassword, name, rInt, name)
}
