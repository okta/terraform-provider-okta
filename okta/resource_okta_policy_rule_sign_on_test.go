package okta

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaPolicyRuleSignon_defaultErrors(t *testing.T) {
	config := testOktaPolicyRuleSignOnDefaultErrors(acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createRuleCheckDestroy(policyRuleSignOn),
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
	oktaIdentityProvider := mgr.GetFixtures("okta_identity_provider.tf", ri, t)
	otherIdentityProvider := mgr.GetFixtures("other_identity_provider.tf", ri, t)
	factorSequence := mgr.GetFixtures("factor_sequence.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", policyRuleSignOn)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createRuleCheckDestroy(policyRuleSignOn),
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
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "access", "DENY"),
					resource.TestCheckResourceAttr(resourceName, "network_connection", "ZONE"),
				),
			},
			{
				Config: oktaIdentityProvider,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "mfa_required", "true"),
				),
			},
			{
				Config: otherIdentityProvider,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "mfa_required", "false"),
					resource.TestCheckResourceAttr(resourceName, "identity_provider", "SPECIFIC_IDP"),
				),
			},

			// This test is failing on our OIE test orgs but not on the non-OIE
			// org. Some orgs need a feature flag for behaviors and/or it isn't
			// supported on OIE orgs
			{
				Config: factorSequence,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "access", "CHALLENGE"),
					resource.TestCheckResourceAttr(resourceName, "behaviors.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "factor_sequence.0.primary_criteria_factor_type", "password"),
					resource.TestCheckResourceAttr(resourceName, "factor_sequence.0.primary_criteria_provider", "OKTA"),
					resource.TestCheckResourceAttr(resourceName, "factor_sequence.0.secondary_criteria.0.factor_type", "push"),
					resource.TestCheckResourceAttr(resourceName, "factor_sequence.0.secondary_criteria.1.factor_type", "token:hotp"),
					resource.TestCheckResourceAttr(resourceName, "factor_sequence.0.secondary_criteria.2.factor_type", "token:software:totp"),
					resource.TestCheckResourceAttr(resourceName, "factor_sequence.0.secondary_criteria.0.provider", "OKTA"),
					resource.TestCheckResourceAttr(resourceName, "factor_sequence.0.secondary_criteria.1.provider", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "factor_sequence.0.secondary_criteria.2.provider", "OKTA"),
					resource.TestCheckResourceAttr(resourceName, "factor_sequence.1.primary_criteria_factor_type", "token:hotp"),
					resource.TestCheckResourceAttr(resourceName, "factor_sequence.1.primary_criteria_provider", "CUSTOM"),
				),
			},
		},
	})
}

func TestAccOktaPolicyRuleSignon_multiple(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(policyRuleSignOn)
	config := mgr.GetFixtures("basic.tf", ri, t)
	basicMultiple := mgr.GetFixtures("basic_multiple.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", policyRuleSignOn)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createRuleCheckDestroy(policyRuleSignOn),
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
				Config: basicMultiple,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(fmt.Sprintf("%s.test_allow", policyRuleSignOn)),
					ensureRuleExists(fmt.Sprintf("%s.test_deny", policyRuleSignOn))),
			},
		},
	})
}

func testOktaPolicyRuleSignOnDefaultErrors(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
	policy_id = "garbageID"
	name     = "Default Rule"
	status   = "ACTIVE"
}
`, policyRuleSignOn, name)
}
