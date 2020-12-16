package okta

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

func deletePolicyRulePasswords(client *testClient) error {
	return deletePolicyRulesByType(sdk.PasswordPolicyType, client)
}

func deletePolicyRulesByType(ruleType string, client *testClient) error {
	ctx := context.Background()
	policies, _, err := client.oktaClient.Policy.ListPolicies(ctx, &query.Params{Type: ruleType})
	if err != nil {
		return fmt.Errorf("failed to list policies in order to properly destroy rules: %v", err)
	}
	for _, policy := range policies {
		rules, _, err := client.apiSupplement.ListPolicyRules(ctx, policy.Id)
		if err != nil {
			return err
		}
		// Tests have always used default policy, I don't really think that is necessarily a good idea but
		// leaving for now, that means we only delete the rules and not the policy, we can keep it around.
		for i := range rules {
			if strings.HasPrefix(rules[i].Name, testResourcePrefix) {
				_, err = client.oktaClient.Policy.DeletePolicyRule(ctx, policy.Id, rules[i].Id)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func TestAccOktaPolicyRulePassword_crud(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicyRulePassword(ri)
	updatedConfig := testOktaPolicyRulePasswordUpdated(ri)
	resourceName := buildResourceFQN(policyRulePassword, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createRuleCheckDestroy(policyRulePassword),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
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
func TestAccOktaPolicyRulePassword_priorityError(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicyRulePriorityError(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createRuleCheckDestroy(policyRulePassword),
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
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createRuleCheckDestroy(policyRulePassword),
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
				return fmt.Errorf("rule still exists, ID: %s, PolicyID: %s", ID, policyID)
			}
		}
		return nil
	}
}

func doesRuleExistsUpstream(policyID, id string) (bool, error) {
	rule, resp, err := getSupplementFromMetadata(testAccProvider.Meta()).GetPolicyRule(context.Background(), policyID, id)
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
	policyid = "${data.okta_default_policy.default-%d.id}"
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
	policyid = "${data.okta_default_policy.default-%d.id}"
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
	policyid = "${data.okta_default_policy.default-%d.id}"
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
	policyid = "${data.okta_default_policy.default-%d.id}"
	name     = "%s"
	status   = "INACTIVE"
	password_change = "DENY"
	password_reset  = "DENY"
	password_unlock = "ALLOW"
}
`, rInt, sdk.PasswordPolicyType, policyRulePassword, name, rInt, name)
}
