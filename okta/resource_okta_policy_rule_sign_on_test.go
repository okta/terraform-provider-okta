package okta

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaPolicyRuleSignon_defaultErrors(t *testing.T) {
	mgr := newFixtureManager("resources", policyRuleSignOn, t.Name())
	config := testOktaPolicyRuleSignOnDefaultErrors(mgr.Seed)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkRuleDestroy(policyRuleSignOn),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Default Rule is immutable"),
			},
		},
	})
}

func TestAccResourceOktaPolicyRuleSignon_crud(t *testing.T) {
	mgr := newFixtureManager("resources", policyRuleSignOn, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	excludedNetwork := mgr.GetFixtures("excluded_network.tf", t)
	oktaIdentityProvider := mgr.GetFixtures("okta_identity_provider.tf", t)
	otherIdentityProvider := mgr.GetFixtures("other_identity_provider.tf", t)
	factorSequence := mgr.GetFixtures("factor_sequence.tf", t)
	resourceName := fmt.Sprintf("%s.test", policyRuleSignOn)

	// NOTE can/will fail with "conditions: Invalid condition type specified: riskScore."
	// Not sure about correct settings for this to pass.
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkRuleDestroy(policyRuleSignOn),
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
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "access", "DENY"),
					resource.TestCheckResourceAttr(resourceName, "network_connection", "ZONE"),
				),
			},
			{
				Config: oktaIdentityProvider,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%d", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "mfa_required", "true"),
				),
			},
			{
				Config: otherIdentityProvider,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%d", mgr.Seed)),
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
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)),
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

func TestAccResourceOktaPolicyRuleSignon_multiple(t *testing.T) {
	mgr := newFixtureManager("resources", policyRuleSignOn, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	basicMultiple := mgr.GetFixtures("basic_multiple.tf", t)
	resourceName := fmt.Sprintf("%s.test", policyRuleSignOn)

	// NOTE can/will fail with "conditions: Invalid condition type specified: riskScore."
	// Not sure about correct settings for this to pass.
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkRuleDestroy(policyRuleSignOn),
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
