package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
)

// TestAccResourceOktaAppSignOnPolicyRules_Issue2797 verifies that setting
// status = "INACTIVE" on a rule is applied correctly and does not cause a
// "provider produced inconsistent result after apply" error.
func TestAccResourceOktaAppSignOnPolicyRules_Issue2797(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppSignOnPolicyRules)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSignOnPolicyRules, t.Name())
	config := mgr.GetFixtures("issue_2797.tf", t)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkAppSignOnPolicyRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "rule.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "rule.1.status", "INACTIVE"),
				),
			},
			{
				// Idempotency check — no diff expected on re-apply.
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

func TestAccResourceOktaAppSignOnPolicyRules_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.policy_rules", resources.OktaIDaaSAppSignOnPolicyRules)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSignOnPolicyRules, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("updated.tf", t)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkAppSignOnPolicyRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "5"),
					// Check rules - rules are stored in config order
					// Rule1: priority 4
					resource.TestCheckResourceAttrSet(resourceName, "rule.0.id"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.name", "Rule1-updatedTF-23/02/26"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.priority", "4"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "rule.0.factor_mode", "2FA"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.network_connection", "ANYWHERE"),
					// Rule2: priority 2
					resource.TestCheckResourceAttr(resourceName, "rule.1.name", "Rule2-updatedTF-23/02/26"),
					resource.TestCheckResourceAttr(resourceName, "rule.1.priority", "2"),
					// Rule3: priority 1
					resource.TestCheckResourceAttr(resourceName, "rule.2.name", "Rule3-updatedTF-23/02/26"),
					resource.TestCheckResourceAttr(resourceName, "rule.2.priority", "1"),
					// Rule4: priority 3
					resource.TestCheckResourceAttr(resourceName, "rule.3.name", "Rule4-updatedTF-23/02/26"),
					resource.TestCheckResourceAttr(resourceName, "rule.3.priority", "3"),
					// Rule5: priority 5
					resource.TestCheckResourceAttr(resourceName, "rule.4.name", "Rule5-updatedTF-23/02/26"),
					resource.TestCheckResourceAttr(resourceName, "rule.4.priority", "5"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "5"),
					// After update - priorities have been shuffled
					// Rule1: priority 4->2
					resource.TestCheckResourceAttr(resourceName, "rule.0.name", "Rule1-updatedTF-23/02/26"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.priority", "2"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.status", idaas.StatusActive),
					// Rule2: priority 2->5
					resource.TestCheckResourceAttr(resourceName, "rule.1.name", "Rule2-updatedTF-23/02/26"),
					resource.TestCheckResourceAttr(resourceName, "rule.1.priority", "5"),
					// Rule3: priority 1->3
					resource.TestCheckResourceAttr(resourceName, "rule.2.name", "Rule3-updatedTF-23/02/26"),
					resource.TestCheckResourceAttr(resourceName, "rule.2.priority", "3"),
					// Rule4: priority 3->1
					resource.TestCheckResourceAttr(resourceName, "rule.3.name", "Rule4-updatedTF-23/02/26"),
					resource.TestCheckResourceAttr(resourceName, "rule.3.priority", "1"),
					// Rule5: priority 5->4
					resource.TestCheckResourceAttr(resourceName, "rule.4.name", "Rule5-updatedTF-23/02/26"),
					resource.TestCheckResourceAttr(resourceName, "rule.4.priority", "4"),
				),
			},
		},
	})
}

// TestAccResourceOktaAppSignOnPolicyRules_dynamicValues tests that the resource
// correctly handles dynamic "rule" blocks where attribute values (e.g. network
// zone IDs) are unknown at plan time. This is a regression test for the
// "Value Conversion Error: Received unknown value, however the target type cannot
// handle unknown values" error that occurred when appSignOnPolicyRulesModel.Rules
// was typed as []policyRuleModel (a native Go slice) instead of types.List.
func TestAccResourceOktaAppSignOnPolicyRules_dynamicValues(t *testing.T) {
	resourceName := fmt.Sprintf("%s.policy_rules", resources.OktaIDaaSAppSignOnPolicyRules)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSignOnPolicyRules, t.Name())
	config := mgr.GetFixtures("dynamic_values.tf", t)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkAppSignOnPolicyRuleDestroy,
		Steps: []resource.TestStep{
			{
				// First apply: network zone ID is unknown at plan time.
				// Previously this triggered "Value Conversion Error" because
				// Rules []policyRuleModel could not hold unknown values.
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
					// 3 rules from the local map
					resource.TestCheckResourceAttr(resourceName, "rule.#", "3"),
					// Rules are stored in config order (map iteration is deterministic
					// in Terraform because it sorts map keys alphabetically).
					// Keys: allow_1fa, allow_2fa, deny_all → sorted: allow_1fa, allow_2fa, deny_all
					resource.TestCheckResourceAttrSet(resourceName, "rule.0.id"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.name", fmt.Sprintf("Allow-1FA-testAcc_%s", mgr.SeedStr())),
					resource.TestCheckResourceAttr(resourceName, "rule.0.factor_mode", "1FA"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.access", "ALLOW"),
					resource.TestCheckResourceAttrSet(resourceName, "rule.1.id"),
					resource.TestCheckResourceAttr(resourceName, "rule.1.name", fmt.Sprintf("Allow-2FA-testAcc_%s", mgr.SeedStr())),
					resource.TestCheckResourceAttr(resourceName, "rule.1.factor_mode", "2FA"),
					resource.TestCheckResourceAttr(resourceName, "rule.1.network_connection", "ZONE"),
					// network_includes should contain the resolved network zone ID
					resource.TestCheckResourceAttrSet(resourceName, "rule.1.network_includes.0"),
					resource.TestCheckResourceAttrSet(resourceName, "rule.2.id"),
					resource.TestCheckResourceAttr(resourceName, "rule.2.name", fmt.Sprintf("Deny-All-testAcc_%s", mgr.SeedStr())),
					resource.TestCheckResourceAttr(resourceName, "rule.2.access", "DENY"),
				),
			},
			{
				// Second apply: idempotency check — no changes expected.
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

// TestAccResourceOktaAppSignOnPolicyRules_Issue2774 tests that LINUX is accepted
// as a valid os_type in platform_include blocks. It was missing from validOSTypes.
func TestAccResourceOktaAppSignOnPolicyRules_Issue2774(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppSignOnPolicyRules)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSignOnPolicyRules, t.Name())
	config := mgr.GetFixtures("issue_2774.tf", t)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkAppSignOnPolicyRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.platform_include.0.os_type", "LINUX"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.platform_include.0.type", "DESKTOP"),
				),
			},
			{
				// Idempotency check — no changes expected after first apply.
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

// TestAccResourceOktaAppSignOnPolicyRules_chains tests that the resource
// correctly handles the chains attribute for authentication chains configuration.
func TestAccResourceOktaAppSignOnPolicyRules_chains(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test_chains", resources.OktaIDaaSAppSignOnPolicyRules)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSignOnPolicyRules, t.Name())
	config := mgr.GetFixtures("chains.tf", t)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkAppSignOnPolicyRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "rule.0.id"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.name", fmt.Sprintf("Chain-Rule-testAcc_%s", mgr.SeedStr())),
					resource.TestCheckResourceAttr(resourceName, "rule.0.chains.#", "2"),
					// Verify both chains are stored
					resource.TestCheckResourceAttrSet(resourceName, "rule.0.chains.0"),
					resource.TestCheckResourceAttrSet(resourceName, "rule.0.chains.1"),
				),
			},
			{
				// Idempotency check
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

// TestAccResourceOktaAppSignOnPolicyRules_keep_me_signed_in verifies that the
// keep_me_signed_in (KMSI / "Option to stay signed in") block on the plural
// resource round-trips correctly across multiple rules. The config defines four
// rules covering every keep_me_signed_in combination:
//   - rule[0]: NOT_ALLOWED with no prompt frequency
//   - rule[1]: ALLOWED with a 50h prompt frequency
//   - rule[2]: ALLOWED with a 168h prompt frequency
//   - rule[3]: NOT_ALLOWED with no prompt frequency (regression case for the
//     "null -> empty string" inconsistent-result-after-apply bug)
//
// It then updates all four rules (flipping post_auth and frequency values) and
// asserts the changes are applied and remain idempotent on re-apply.
func TestAccResourceOktaAppSignOnPolicyRules_keep_me_signed_in(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppSignOnPolicyRules)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSignOnPolicyRules, t.Name())
	config := mgr.GetFixtures("keep_me_signed_in.tf", t)
	updatedConfig := mgr.GetFixtures("keep_me_signed_in_updated.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkAppSignOnPolicyRuleDestroy,
		Steps: []resource.TestStep{
			{
				// Step 1: Create four rules with varying KMSI settings.
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "4"),
					// rule[0]: NOT_ALLOWED, no frequency. The Optional (non-Computed)
					// frequency is null, so it is absent from state (not "").
					resource.TestCheckResourceAttr(resourceName, "rule.0.access", "ALLOW"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.keep_me_signed_in.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.keep_me_signed_in.0.post_auth", "NOT_ALLOWED"),
					resource.TestCheckNoResourceAttr(resourceName, "rule.0.keep_me_signed_in.0.post_auth_prompt_frequency"),
					// rule[1]: ALLOWED, PT50H.
					resource.TestCheckResourceAttr(resourceName, "rule.1.keep_me_signed_in.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.1.keep_me_signed_in.0.post_auth", "ALLOWED"),
					resource.TestCheckResourceAttr(resourceName, "rule.1.keep_me_signed_in.0.post_auth_prompt_frequency", "PT50H"),
					// rule[2]: ALLOWED, PT168H.
					resource.TestCheckResourceAttr(resourceName, "rule.2.keep_me_signed_in.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.2.keep_me_signed_in.0.post_auth", "ALLOWED"),
					resource.TestCheckResourceAttr(resourceName, "rule.2.keep_me_signed_in.0.post_auth_prompt_frequency", "PT168H"),
					// rule[3]: NOT_ALLOWED, no frequency. This is the exact scenario
					// from the bug report (rule[3].keep_me_signed_in[0].
					// post_auth_prompt_frequency was null, but now ""): it must stay
					// null after apply, i.e. absent from state.
					resource.TestCheckResourceAttr(resourceName, "rule.3.keep_me_signed_in.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.3.keep_me_signed_in.0.post_auth", "NOT_ALLOWED"),
					resource.TestCheckNoResourceAttr(resourceName, "rule.3.keep_me_signed_in.0.post_auth_prompt_frequency"),
				),
			},
			{
				// Step 2: Idempotency on create config.
				Config:   config,
				PlanOnly: true,
			},
			{
				// Step 3: Update all four rules, flipping post_auth and frequency.
				// rule[1] clears its frequency (ALLOWED PT50H -> NOT_ALLOWED), which
				// is the regression scenario: the API returns an empty frequency and
				// the provider must keep the Optional (non-Computed) attribute null
				// instead of "" to avoid an "inconsistent result after apply" error.
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "rule.#", "4"),
					// rule[0]: NOT_ALLOWED -> ALLOWED (PT168H).
					resource.TestCheckResourceAttr(resourceName, "rule.0.keep_me_signed_in.0.post_auth", "ALLOWED"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.keep_me_signed_in.0.post_auth_prompt_frequency", "PT168H"),
					// rule[1]: ALLOWED (PT50H) -> NOT_ALLOWED, frequency cleared
					// (null, so absent from state).
					resource.TestCheckResourceAttr(resourceName, "rule.1.keep_me_signed_in.0.post_auth", "NOT_ALLOWED"),
					resource.TestCheckNoResourceAttr(resourceName, "rule.1.keep_me_signed_in.0.post_auth_prompt_frequency"),
					// rule[2]: ALLOWED (PT168H) -> ALLOWED (PT50H).
					resource.TestCheckResourceAttr(resourceName, "rule.2.keep_me_signed_in.0.post_auth", "ALLOWED"),
					resource.TestCheckResourceAttr(resourceName, "rule.2.keep_me_signed_in.0.post_auth_prompt_frequency", "PT50H"),
					// rule[3]: NOT_ALLOWED -> ALLOWED (PT168H).
					resource.TestCheckResourceAttr(resourceName, "rule.3.keep_me_signed_in.0.post_auth", "ALLOWED"),
					resource.TestCheckResourceAttr(resourceName, "rule.3.keep_me_signed_in.0.post_auth_prompt_frequency", "PT168H"),
				),
			},
			{
				// Step 4: Idempotency on updated config.
				Config:   updatedConfig,
				PlanOnly: true,
			},
		},
	})
}
