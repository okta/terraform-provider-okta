package idaas_test

import (
	"context"
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

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkPostAuthSessionPolicyRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
					resource.TestCheckResourceAttr(resourceName, "terminate_session", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: mgr.ConfigReplace(updatedConfig),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
					resource.TestCheckResourceAttr(resourceName, "terminate_session", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("not found: %s", resourceName)
					}
					return fmt.Sprintf("%s/%s", rs.Primary.Attributes["policy_id"], rs.Primary.ID), nil
				},
			},
		},
	})
}

func checkPostAuthSessionPolicyRuleDestroy(s *terraform.State) error {
	client := iDaaSAPIClientForTestUtil.OktaSDKClientV6()
	for _, r := range s.RootModule().Resources {
		if r.Type != resources.OktaIDaaSPostAuthSessionPolicyRule {
			continue
		}
		policyId := r.Primary.Attributes["policy_id"]
		ruleId := r.Primary.ID

		resp, _, err := client.PolicyAPI.GetPolicyRule(context.Background(), policyId, ruleId).Execute()
		if err == nil && resp.PostAuthSessionPolicyRule != nil {
			return fmt.Errorf("post auth session policy rule still exists")
		}
	}
	return nil
}
