package okta

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	articulateOkta "github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/okta/okta-sdk-golang/okta"
)

func deletePasswordPolicyRules(artClient *articulateOkta.Client, client *okta.Client) error {
	return deletePolicyRulesByType(passwordPolicyType, artClient, client)
}

func deletePolicyRulesByType(ruleType string, artClient *articulateOkta.Client, client *okta.Client) error {
	policies, _, err := artClient.Policies.GetPoliciesByType(ruleType)

	if err != nil {
		return err
	}

	for _, policy := range policies.Policies {
		rules, _, err := artClient.Policies.GetPolicyRules(policy.ID)

		if err == nil {
			// Tests have always used default policy, I don't really think that is necessarily a good idea but
			// leaving for now, that means we only delete the rules and not the policy, we can keep it around.
			for _, rule := range rules.Rules {
				if strings.HasPrefix(rule.Name, testResourcePrefix) {
					_, err = artClient.Policies.DeletePolicyRule(policy.ID, rule.ID)
				}
			}
		}
	}

	return err
}

func TestAccOktaPolicyRulePassword(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicyRulePassword(ri)
	updatedConfig := testOktaPolicyRulePassword_updated(ri)
	resourceName := buildResourceFQN(passwordPolicyRule, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOktaPolicyRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyRuleExists(resourceName),
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
func TestAccOktaPolicyRulePassword_signonErrors(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicyRulePassword(ri)
	updatedConfig := testOktaPolicyRulePassword_signonErrors(ri)
	resourceName := buildResourceFQN(passwordPolicyRule, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOktaPolicyRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyRuleExists(resourceName),
				),
			},
			{
				Config:      updatedConfig,
				ExpectError: regexp.MustCompile("config is invalid: .* invalid or unknown key: session_idle"),
				PlanOnly:    true,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyRuleExists(resourceName),
				),
			},
		},
	})
}
func TestAccOktaPolicyRulePassword_authErrors(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicyRulePassword(ri)
	updatedConfig := testOktaPolicyRulePassword_authtErrors(ri)
	resourceName := buildResourceFQN(passwordPolicyRule, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOktaPolicyRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyRuleExists(resourceName),
				),
			},
			{
				Config:      updatedConfig,
				ExpectError: regexp.MustCompile("config is invalid: .* : invalid or unknown key: auth_type"),
				PlanOnly:    true,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyRuleExists(resourceName),
				),
			},
		},
	})
}

func testOktaPolicyRuleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("[ERROR] Resource Not found: %s", name)
		}

		policyID, hasP := rs.Primary.Attributes["policyid"]
		if !hasP {
			return fmt.Errorf("[ERROR] No policy ID found in state for Policy Rule")
		}
		ruleID, hasR := rs.Primary.Attributes["id"]
		if !hasR {
			return fmt.Errorf("[ERROR] No rule ID found in state for Policy Rule")
		}
		ruleName, hasName := rs.Primary.Attributes["name"]
		if !hasName {
			return fmt.Errorf("[ERROR] No name found in state for Policy Rule")
		}

		err := testPolicyRuleExists(true, policyID, ruleID, ruleName)
		if err != nil {
			return err
		}
		return nil
	}
}

func testOktaPolicyRuleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "okta_policy_rules" {
			continue
		}

		policyID, hasP := rs.Primary.Attributes["policyid"]
		if !hasP {
			return fmt.Errorf("[ERROR] No policy ID found in state for Policy Rule")
		}
		ruleID, hasR := rs.Primary.Attributes["id"]
		if !hasR {
			return fmt.Errorf("[ERROR] No rule ID found in state for Policy Rule")
		}
		ruleName, hasName := rs.Primary.Attributes["name"]
		if !hasName {
			return fmt.Errorf("[ERROR] No name found in state for Policy Rule")
		}

		err := testPolicyRuleExists(false, policyID, ruleID, ruleName)
		if err != nil {
			return err
		}
	}
	return nil
}

func testPolicyRuleExists(expected bool, policyID string, ruleID, ruleName string) error {
	client := testAccProvider.Meta().(*Config).articulateOktaClient

	exists := false
	_, _, err := client.Policies.GetPolicy(policyID)
	if err != nil {
		if client.OktaErrorCode != "E0000007" {
			return fmt.Errorf("[ERROR] Error Listing Policy in Okta: %v", err)
		}
	} else {
		_, _, err := client.Policies.GetPolicyRule(policyID, ruleID)
		if err != nil {
			if client.OktaErrorCode != "E0000007" {
				return fmt.Errorf("[ERROR] Error Listing Policy Rule in Okta: %v", err)
			}
		} else {
			exists = true
		}
	}

	if expected == true && exists == false {
		return fmt.Errorf("[ERROR] Policy Rule %v not found in Okta", ruleName)
	} else if expected == false && exists == true {
		return fmt.Errorf("[ERROR] Policy Rule %v still exists in Okta", ruleName)
	}
	return nil
}

func testOktaPolicyRulePassword(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
data "okta_default_policies" "default-%d" {
	type = "%s"
}

resource "%s" "%s" {
	policyid = "${data.okta_default_policies.default-%d.id}"
	name     = "%s"
	status   = "ACTIVE"
}
`, rInt, passwordPolicyType, passwordPolicyRule, name, rInt, name)
}

func testOktaPolicyRulePassword_updated(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
data "okta_default_policies" "default-%d" {
	type = "%s"
}

resource "%s" "%s" {
	policyid = "${data.okta_default_policies.default-%d.id}"
	name     = "%s"
	status   = "INACTIVE"
	password_change = "DENY"
	password_reset  = "DENY"
	password_unlock = "ALLOW"
}
`, rInt, passwordPolicyType, passwordPolicyRule, name, rInt, name)
}

func testOktaPolicyRulePassword_signonErrors(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
data "okta_default_policies" "default-%d" {
	type = "%s"
}

resource "%s" "%s" {
	policyid = "${data.okta_default_policies.default-%d.id}"
	name     = "%s"
	status   = "ACTIVE"
	session_idle = 240
}
`, rInt, passwordPolicyType, passwordPolicyRule, name, rInt, name)
}

func testOktaPolicyRulePassword_authtErrors(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
data "okta_default_policies" "default-%d" {
	type = "%s"
}

resource "%s" "%s" {
	policyid = "${data.okta_default_policies.default-%d.id}"
	name     = "%s"
	status   = "ACTIVE"
	auth_type = "RADIUS"
}
`, rInt, passwordPolicyType, passwordPolicyRule, name, rInt, name)
}
