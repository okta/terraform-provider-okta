package okta

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func deleteSignOnPolicyRules(client *testClient) error {
	return deletePolicyRulesByType(signOnPolicyType, client)
}

func TestAccOktaPolicyRuleDefaultErrors(t *testing.T) {
	config := testOktaPolicyRuleSignOnDefaultErrors(acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createRuleCheckDestroy(policyRuleSignOn),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Default Rule is immutable"),
			},
		},
	})
}

func TestAccOktaPolicyRuleSignOn(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(policyRuleSignOn)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", policyRuleSignOn)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createRuleCheckDestroy(policyRuleSignOn),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "access", "DENY"),
					resource.TestCheckResourceAttr(resourceName, "session_idle", "240"),
					resource.TestCheckResourceAttr(resourceName, "session_lifetime", "240"),
					resource.TestCheckResourceAttr(resourceName, "session_persistent", "false"),
					resource.TestCheckResourceAttr(resourceName, "users_excluded.#", "1"),
				),
			},
		},
	})
}

func TestAccOktaPolicyRuleSignOnPassErrors(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicyRuleSignOnPassErrors(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createRuleCheckDestroy(policyRuleSignOn),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("config is invalid: .*: : invalid or unknown key: password_change"),
				PlanOnly:    true,
			},
		},
	})
}

func testOktaPolicyRuleSignOnDefaultErrors(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
	policyid = "garbageID"
	name     = "Default Rule"
	status   = "ACTIVE"
}
`, policyRuleSignOn, name)
}

func testOktaPolicyRuleSignOnPassErrors(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
data "okta_default_policy" "default-%d" {
	type = "%s"
}

resource "%s" "%s" {
  policyid = "${data.okta_default_policy.default-%d.id}"
  name     = "%s"
  status   = "ACTIVE"
  password_change = "DENY"
}
`, rInt, signOnPolicyType, policyRuleSignOn, name, rInt, name)
}
