package idaas_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
)

func TestAccResourceOktaPolicyRuleSignon_defaultErrors(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSPolicyRuleSignOn, t.Name())
	config := testOktaPolicyRuleSignOnDefaultErrors(mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkRuleDestroy(resources.OktaIDaaSPolicyRuleSignOn),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Default Rule is immutable"),
			},
		},
	})
}

func TestAccResourceOktaPolicyRuleSignon_GH2419(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSPolicyRuleSignOn, t.Name())
	config := mgr.GetFixtures("gh2419.tf", t)
	updatedConfig := mgr.GetFixtures("gh2419_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test_risk_ONLY", resources.OktaIDaaSPolicyRuleSignOn)
	resourceName2 := fmt.Sprintf("%s.test_risc_ONLY", resources.OktaIDaaSPolicyRuleSignOn)
	resourceName3 := fmt.Sprintf("%s.test_BOTH", resources.OktaIDaaSPolicyRuleSignOn)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkRuleDestroy(resources.OktaIDaaSPolicyRuleSignOn),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName2, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName3, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "risk_level", "ANY"),
					resource.TestCheckResourceAttr(resourceName2, "risc_level", "MEDIUM"),
					resource.TestCheckResourceAttr(resourceName3, "risk_level", "LOW"),
					resource.TestCheckResourceAttr(resourceName3, "risc_level", "HIGH"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName2, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName3, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "risk_level", "MEDIUM"),
					resource.TestCheckResourceAttr(resourceName2, "risc_level", "HIGH"),
					resource.TestCheckResourceAttr(resourceName3, "risk_level", "MEDIUM"),
					resource.TestCheckResourceAttr(resourceName3, "risc_level", "HIGH"),
				),
			},
		},
	})
}

func TestAccResourceOktaPolicyRuleSignon_GH2494(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSPolicyRuleSignOn, t.Name())
	config := mgr.GetFixtures("gh2494.tf", t)
	updatedConfig := mgr.GetFixtures("gh2494_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test_NEITHER", resources.OktaIDaaSPolicyRuleSignOn)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkRuleDestroy(resources.OktaIDaaSPolicyRuleSignOn),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckNoResourceAttr(resourceName, "risc_level"),
					resource.TestCheckNoResourceAttr(resourceName, "risk_level"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckNoResourceAttr(resourceName, "risk_level"),
					resource.TestCheckNoResourceAttr(resourceName, "risc_level"),
				),
			},
		},
	})
}

func TestAccResourceOktaPolicyRuleSignon_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSPolicyRuleSignOn, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	excludedNetwork := mgr.GetFixtures("excluded_network.tf", t)
	oktaIdentityProvider := mgr.GetFixtures("okta_identity_provider.tf", t)
	otherIdentityProvider := mgr.GetFixtures("other_identity_provider.tf", t)
	factorSequence := mgr.GetFixtures("factor_sequence.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSPolicyRuleSignOn)

	// NOTE can/will fail with "conditions: Invalid condition type specified: riskScore."
	// Not sure about correct settings for this to pass.
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkRuleDestroy(resources.OktaIDaaSPolicyRuleSignOn),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusInactive),
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
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "access", "DENY"),
					resource.TestCheckResourceAttr(resourceName, "network_connection", "ZONE"),
				),
			},
			{
				Config: oktaIdentityProvider,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%d", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "mfa_required", "true"),
				),
			},
			{
				Config: otherIdentityProvider,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%d", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
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
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
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
	mgr := newFixtureManager("resources", resources.OktaIDaaSPolicyRuleSignOn, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	basicMultiple := mgr.GetFixtures("basic_multiple.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSPolicyRuleSignOn)

	// NOTE can/will fail with "conditions: Invalid condition type specified: riskScore."
	// Not sure about correct settings for this to pass.
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkRuleDestroy(resources.OktaIDaaSPolicyRuleSignOn),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
				),
			},
			{
				Config: basicMultiple,
				Check: resource.ComposeTestCheckFunc(
					ensureRuleExists(fmt.Sprintf("%s.test_allow", resources.OktaIDaaSPolicyRuleSignOn)),
					ensureRuleExists(fmt.Sprintf("%s.test_deny", resources.OktaIDaaSPolicyRuleSignOn))),
			},
		},
	})
}

func testOktaPolicyRuleSignOnDefaultErrors(rInt int) string {
	name := acctest.BuildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
	policy_id = "garbageID"
	name     = "Default Rule"
	status   = "ACTIVE"
}
`, resources.OktaIDaaSPolicyRuleSignOn, name)
}
