package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
)

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
