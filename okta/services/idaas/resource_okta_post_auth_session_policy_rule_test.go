package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaPostAuthSessionPolicyRule_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSPostAuthSessionPolicyRule)
	mgr := newFixtureManager("resources", resources.OktaIDaaSPostAuthSessionPolicyRule, t.Name())

	dataSourceConfig := `
	data "okta_post_auth_session_policy" "test" {
	}
	`

	config := `
	data "okta_post_auth_session_policy" "test" {
	}

	resource "okta_post_auth_session_policy_rule" "test" {
	policy_id         = data.okta_post_auth_session_policy.test.id
	name              = "testAcc-replace_with_uuid"
	terminate_session = true
	}
	`

	updatedConfig := `
	data "okta_post_auth_session_policy" "test" {
	}

	resource "okta_post_auth_session_policy_rule" "test" {
	policy_id         = data.okta_post_auth_session_policy.test.id
	name              = "testAcc-replace_with_uuid"
	terminate_session = false
	}
	`

	var policyId string
	var ruleId string

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				// Step 1: Apply data source only to get policy ID and rule ID
				Config: mgr.ConfigReplace(dataSourceConfig),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_post_auth_session_policy.test", "id"),
					// Capture the policy ID and rule ID for import
					func(s *terraform.State) error {
						ds, ok := s.RootModule().Resources["data.okta_post_auth_session_policy.test"]
						if !ok {
							return fmt.Errorf("data source not found")
						}
						policyId = ds.Primary.ID
						ruleId = ds.Primary.Attributes["rule_id"]
						return nil
					},
				),
			},
			{
				// Step 2: Import the existing rule
				Config:             mgr.ConfigReplace(config),
				ResourceName:       resourceName,
				ImportState:        true,
				ImportStatePersist: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					return fmt.Sprintf("%s/%s", policyId, ruleId), nil
				},
			},
			{
				// Step 3: Apply config to update rule
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
					resource.TestCheckResourceAttr(resourceName, "terminate_session", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				// Step 4: Update
				Config: mgr.ConfigReplace(updatedConfig),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
					resource.TestCheckResourceAttr(resourceName, "terminate_session", "false"),
				),
			},
		},
	})
}
