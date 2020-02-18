package okta

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func deleteSignOnPolicyRules(client *testClient) error {
	return deletePolicyRulesByType(signOnPolicyType, client)
}

func TestAccOktaPolicyRuleSignon_defaultErrors(t *testing.T) {
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

func TestAccOktaPolicyRuleSignon_crud(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(policyRuleSignOn)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", ri, t)
	excludedNetwork := mgr.GetFixtures("excluded_network.tf", ri, t)
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
			{
				Config: excludedNetwork,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "access", "DENY"),
					resource.TestCheckResourceAttr(resourceName, "network_connection", "ZONE"),
				),
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
