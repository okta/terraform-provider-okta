package okta

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func deletepolicyRulePasswords(client *testClient) error {
	return deletePolicyRulesByType(passwordPolicyType, client)
}

func deletePolicyRulesByType(ruleType string, client *testClient) error {
	policies, _, err := client.artClient.Policies.GetPoliciesByType(ruleType)

	if err != nil {
		return err
	}

	for _, policy := range policies.Policies {
		rules, _, err := client.artClient.Policies.GetPolicyRules(policy.ID)

		if err == nil && rules != nil {
			// Tests have always used default policy, I don't really think that is necessarily a good idea but
			// leaving for now, that means we only delete the rules and not the policy, we can keep it around.
			for _, rule := range rules.Rules {
				if strings.HasPrefix(rule.Name, testResourcePrefix) {
					_, err = client.artClient.Policies.DeletePolicyRule(policy.ID, rule.ID)
				}
			}
		}
	}

	return err
}

func TestAccOktaPolicyRulePassword_crud(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicyRulePassword(ri)
	updatedConfig := testOktaPolicyRulePasswordUpdated(ri)
	resourceName := buildResourceFQN(policyRulePassword, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createRuleCheckDestroy(policyRulePassword),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "password_change", "DENY"),
					resource.TestCheckResourceAttr(resourceName, "password_reset", "DENY"),
					resource.TestCheckResourceAttr(resourceName, "password_unlock", "ALLOW"),
				),
			},
		},
	})
}

// Testing the logic that errors when an invalid priority is provided
func TestAccOktaPolicyRulePassword_priorityError(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicyRulePriorityError(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createRuleCheckDestroy(policyRulePassword),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("provided priority was not valid, got: 999, API responded with: 1. See schema for attribute details"),
			},
		},
	})
}

// Testing the successful setting of priority
func TestAccOktaPolicyRulePassword_priority(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicyRulePriority(ri)
	resourceName := buildResourceFQN(policyRulePassword, ri)
	name := buildResourceName(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createRuleCheckDestroy(policyRulePassword),
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

func ensureRuleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		missingErr := fmt.Errorf("resource not found: %s", name)
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return missingErr
		}

		policyID := rs.Primary.Attributes["policyid"]
		ID := rs.Primary.ID
		exist, err := doesRuleExistsUpstream(policyID, ID)
		if err != nil {
			return err
		} else if !exist {
			return missingErr
		}

		return nil
	}
}

func createRuleCheckDestroy(ruleType string) func(*terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != ruleType {
				continue
			}

			policyID := rs.Primary.Attributes["policyid"]
			ID := rs.Primary.ID
			exists, err := doesRuleExistsUpstream(policyID, ID)
			if err != nil {
				return err
			}

			if exists {
				return fmt.Errorf("Rule still exists, ID: %s, PolicyID: %s", ID, policyID)
			}
		}
		return nil
	}
}

func doesRuleExistsUpstream(policyID string, ID string) (bool, error) {
	client := getClientFromMetadata(testAccProvider.Meta())

	rule, resp, err := client.Policies.GetPolicyRule(policyID, ID)
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return rule.ID != "", nil
}

func testOktaPolicyRulePassword(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
data "okta_default_policy" "default-%d" {
	type = "%s"
}

resource "%s" "%s" {
	policyid = "${data.okta_default_policy.default-%d.id}"
	name     = "%s"
	status   = "ACTIVE"
}
`, rInt, passwordPolicyType, policyRulePassword, name, rInt, name)
}

func testOktaPolicyRulePriority(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
data "okta_default_policy" "default-%d" {
	type = "%s"
}

resource "%s" "%s" {
	policyid = "${data.okta_default_policy.default-%d.id}"
	name     = "%s"
	priority = 1
	status   = "ACTIVE"
}
`, rInt, passwordPolicyType, policyRulePassword, name, rInt, name)
}

func testOktaPolicyRulePriorityError(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
data "okta_default_policy" "default-%d" {
	type = "%s"
}

resource "%s" "%s" {
	policyid = "${data.okta_default_policy.default-%d.id}"
	name     = "%s"
	priority = 999
	status   = "ACTIVE"
}
`, rInt, passwordPolicyType, policyRulePassword, name, rInt, name)
}

func testOktaPolicyRulePasswordUpdated(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
data "okta_default_policy" "default-%d" {
	type = "%s"
}

resource "%s" "%s" {
	policyid = "${data.okta_default_policy.default-%d.id}"
	name     = "%s"
	status   = "INACTIVE"
	password_change = "DENY"
	password_reset  = "DENY"
	password_unlock = "ALLOW"
}
`, rInt, passwordPolicyType, policyRulePassword, name, rInt, name)
}

func testOktaPolicyRulePasswordSignOnErrors(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
data "okta_default_policy" "default-%d" {
	type = "%s"
}

resource "%s" "%s" {
	policyid = "${data.okta_default_policy.default-%d.id}"
	name     = "%s"
	status   = "ACTIVE"
	session_idle = 240
}
`, rInt, passwordPolicyType, policyRulePassword, name, rInt, name)
}

func testOktaPolicyRulePasswordAuthErrors(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
data "okta_default_policy" "default-%d" {
	type = "%s"
}

resource "%s" "%s" {
	policyid = "${data.okta_default_policy.default-%d.id}"
	name     = "%s"
	status   = "ACTIVE"
	auth_type = "RADIUS"
}
`, rInt, passwordPolicyType, policyRulePassword, name, rInt, name)
}
