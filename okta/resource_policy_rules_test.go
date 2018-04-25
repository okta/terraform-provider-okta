package okta

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOktaPolicyRule_defaultErrors(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicyRuleSignOn_defaultErrors(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("You cannot edit a default Policy Rule"),
				PlanOnly:    true,
			},
		},
	})
}

func TestAccOktaPolicyRules_nameErrors(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicyRuleSignOn(ri)
	updatedConfig := testOktaPolicyRuleSignOn_nameErrors(ri)
	resourceName := "okta_policy_rules.test-" + strconv.Itoa(ri)

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
				ExpectError: regexp.MustCompile("You cannot change the policyid field, name field, or type field of an existing Policy Rule"),
				PlanOnly:    true,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyRuleExists(resourceName),
				),
			},
		},
	})
}
func TestAccOktaPolicyRules_typeErrors(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicyRuleSignOn(ri)
	updatedConfig := testOktaPolicyRuleSignOn_typeErrors(ri)
	resourceName := "okta_policy_rules.test-" + strconv.Itoa(ri)

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
				ExpectError: regexp.MustCompile("You cannot change the policyid field, name field, or type field of an existing Policy Rule"),
				PlanOnly:    true,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyRuleExists(resourceName),
				),
			},
		},
	})
}
func TestAccOktaPolicyRules_policyErrors(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicyRuleSignOn(ri)
	updatedConfig := testOktaPolicyRuleSignOn_policyErrors(ri)
	resourceName := "okta_policy_rules.test-" + strconv.Itoa(ri)

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
				ExpectError: regexp.MustCompile("You cannot change the policyid field, name field, or type field of an existing Policy Rule"),
				PlanOnly:    true,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyRuleExists(resourceName),
				),
			},
		},
	})
}
func TestAccOktaPolicyRuleSignOn(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicyRuleSignOn(ri)
	updatedConfig := testOktaPolicyRuleSignOn_updated(ri)
	resourceName := "okta_policy_rules.test-" + strconv.Itoa(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOktaPolicyRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "SIGN_ON"),
					resource.TestCheckResourceAttr(resourceName, "name", "testAcc-"+strconv.Itoa(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "SIGN_ON"),
					resource.TestCheckResourceAttr(resourceName, "name", "testAcc-"+strconv.Itoa(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "priority", "999"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.signon.0.access", "DENY"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.signon.0.sessionidle", "240"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.signon.0.sessionlifetime", "240"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.signon.0.persistentcookie", "false"),
				),
			},
		},
	})
}

func TestAccOktaPolicyRuleSignOn_passErrors(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicyRuleSignOn(ri)
	updatedConfig := testOktaPolicyRuleSignOn_passErrors(ri)
	resourceName := "okta_policy_rules.test-" + strconv.Itoa(ri)

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
				ExpectError: regexp.MustCompile("password action options not supported in the Okta SignOn Policy Rule"),
				PlanOnly:    true,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyRuleExists(resourceName),
				),
			},
		},
	})
}

func TestAccOktaPolicyRulePassword(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicyRulePassword(ri)
	updatedConfig := testOktaPolicyRulePassword_updated(ri)
	resourceName := "okta_policy_rules.test-" + strconv.Itoa(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOktaPolicyRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "PASSWORD"),
					resource.TestCheckResourceAttr(resourceName, "name", "testAcc-"+strconv.Itoa(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "PASSWORD"),
					resource.TestCheckResourceAttr(resourceName, "name", "testAcc-"+strconv.Itoa(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "priority", "999"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.password.0.passwordchange", "DENY"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.password.0.passwordreset", "DENY"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.password.0.passwordunlock", "ALLOW"),
				),
			},
		},
	})
}
func TestAccOktaPolicyRulePassword_signonErrors(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicyRulePassword(ri)
	updatedConfig := testOktaPolicyRulePassword_signonErrors(ri)
	resourceName := "okta_policy_rules.test-" + strconv.Itoa(ri)

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
				ExpectError: regexp.MustCompile("signon action options not supported in the Okta Password Policy Rule"),
				PlanOnly:    true,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyRuleExists(resourceName),
				),
			},
		},
	})
}
func TestAccOktaPolicyRulePassword_authtErrors(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicyRulePassword(ri)
	updatedConfig := testOktaPolicyRulePassword_authtErrors(ri)
	resourceName := "okta_policy_rules.test-" + strconv.Itoa(ri)

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
				ExpectError: regexp.MustCompile("authtype condition options not supported in the Okta Password Policy Rule"),
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
	return nil
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
	client := testAccProvider.Meta().(*Config).oktaClient

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

func testOktaPolicyRuleSignOn(rInt int) string {
	return fmt.Sprintf(`
data "okta_default_policies" "default-%d" {
  type = "OKTA_SIGN_ON"
}

resource "okta_policy_rules" "test-%d" {
  policyid = "${data.okta_default_policies.default-%d.id}"
  type     = "SIGN_ON"
  name     = "testAcc-%d"
  status   = "ACTIVE"
}
`, rInt, rInt, rInt, rInt)
}

func testOktaPolicyRuleSignOn_updated(rInt int) string {
	return fmt.Sprintf(`
data "okta_default_policies" "default-%d" {
  type = "OKTA_SIGN_ON"
}

resource "okta_policy_rules" "test-%d" {
  policyid = "${data.okta_default_policies.default-%d.id}"
  type     = "SIGN_ON"
  name     = "testAcc-%d"
  priority = 999
  status   = "INACTIVE"
  actions {
    signon {
      access           = "DENY"
      sessionidle      = 240
      sessionlifetime  = 240
      persistentcookie = false
    }
  }
}
`, rInt, rInt, rInt, rInt)
}

func testOktaPolicyRuleSignOn_defaultErrors(rInt int) string {
	return fmt.Sprintf(`
resource "okta_policy_rules" "test-%d" {
  policyid = "garbageID"
  type     = "SIGN_ON"
  name     = "Default Rule"
  status   = "ACTIVE"
}
`, rInt)
}

func testOktaPolicyRuleSignOn_nameErrors(rInt int) string {
	return fmt.Sprintf(`
data "okta_default_policies" "default-%d" {
  type = "OKTA_SIGN_ON"
}

resource "okta_policy_rules" "test-%d" {
  policyid = "${data.okta_default_policies.default-%d.id}"
  type     = "SIGN_ON"
  name     = "testAccChanged-%d"
  status   = "ACTIVE"
}
`, rInt, rInt, rInt, rInt)
}

func testOktaPolicyRuleSignOn_typeErrors(rInt int) string {
	return fmt.Sprintf(`
data "okta_default_policies" "default-%d" {
  type = "OKTA_SIGN_ON"
}

resource "okta_policy_rules" "test-%d" {
  policyid = "${data.okta_default_policies.default-%d.id}"
  type     = "PASSWORD"
  name     = "testAcc-%d"
  status   = "ACTIVE"
}
`, rInt, rInt, rInt, rInt)
}

func testOktaPolicyRuleSignOn_policyErrors(rInt int) string {
	return fmt.Sprintf(`
data "okta_default_policies" "default-%d" {
  type = "OKTA_SIGN_ON"
}

resource "okta_policy_rules" "test-%d" {
  policyid = "changedPolicyID"
  type     = "SIGN_ON"
  name     = "testAcc-%d"
  status   = "ACTIVE"
}
`, rInt, rInt, rInt)
}

func testOktaPolicyRuleSignOn_passErrors(rInt int) string {
	return fmt.Sprintf(`
data "okta_default_policies" "default-%d" {
  type = "OKTA_SIGN_ON"
}

resource "okta_policy_rules" "test-%d" {
  policyid = "${data.okta_default_policies.default-%d.id}"
  type     = "SIGN_ON"
  name     = "testAcc-%d"
  status   = "ACTIVE"
  actions {
    password {
      passwordchange = "DENY"
    }
  }
}
`, rInt, rInt, rInt, rInt)
}

func testOktaPolicyRulePassword(rInt int) string {
	return fmt.Sprintf(`
data "okta_default_policies" "default-%d" {
  type = "PASSWORD"
}

resource "okta_policy_rules" "test-%d" {
  policyid = "${data.okta_default_policies.default-%d.id}"
  type     = "PASSWORD"
  name     = "testAcc-%d"
  status   = "ACTIVE"
}
`, rInt, rInt, rInt, rInt)
}

func testOktaPolicyRulePassword_updated(rInt int) string {
	return fmt.Sprintf(`
data "okta_default_policies" "default-%d" {
  type = "PASSWORD"
}

resource "okta_policy_rules" "test-%d" {
  policyid = "${data.okta_default_policies.default-%d.id}"
  type     = "PASSWORD"
  name     = "testAcc-%d"
  priority = 999
  status   = "INACTIVE"
  actions {
    password {
      passwordchange = "DENY"
      passwordreset  = "DENY"
      passwordunlock = "ALLOW"
    }
  }
}
`, rInt, rInt, rInt, rInt)
}

func testOktaPolicyRulePassword_signonErrors(rInt int) string {
	return fmt.Sprintf(`
data "okta_default_policies" "default-%d" {
  type = "PASSWORD"
}

resource "okta_policy_rules" "test-%d" {
  policyid = "${data.okta_default_policies.default-%d.id}"
  type     = "PASSWORD"
  name     = "testAcc-%d"
  status   = "ACTIVE"
  actions {
    signon {
      sessionidle = 240
    }
  }
}
`, rInt, rInt, rInt, rInt)
}

func testOktaPolicyRulePassword_authtErrors(rInt int) string {
	return fmt.Sprintf(`
data "okta_default_policies" "default-%d" {
  type = "PASSWORD"
}

resource "okta_policy_rules" "test-%d" {
  policyid = "${data.okta_default_policies.default-%d.id}"
  type     = "PASSWORD"
  name     = "testAcc-%d"
  status   = "ACTIVE"
  conditions {
    authtype = "RADIUS"
  }
}
`, rInt, rInt, rInt, rInt)
}
