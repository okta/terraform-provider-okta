package idaas_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
	"testing"
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
