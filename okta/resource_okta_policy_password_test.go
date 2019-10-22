package okta

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func deletePasswordPolicies(client *testClient) error {
	return deletePolicyByType(passwordPolicyType, client)
}

func deletePolicyByType(t string, client *testClient) error {
	col, _, err := client.artClient.Policies.GetPoliciesByType(t)

	if err != nil {
		return fmt.Errorf("Failed to retrieve policies in order to properly destroy. Error: %s", err)
	}

	for _, policy := range col.Policies {
		if strings.HasPrefix(policy.Name, testResourcePrefix) {
			_, err = client.artClient.Policies.DeletePolicy(policy.ID)
		}
	}

	return nil
}

func TestAccOktaPolicyPassword_crud(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(policyPassword)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", policyPassword)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createPolicyCheckDestroy(policyPassword),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test Password Policy"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test Password Policy Updated"),
					resource.TestCheckResourceAttr(resourceName, "password_min_length", "12"),
					resource.TestCheckResourceAttr(resourceName, "password_min_lowercase", "0"),
					resource.TestCheckResourceAttr(resourceName, "password_min_uppercase", "0"),
					resource.TestCheckResourceAttr(resourceName, "password_min_number", "0"),
					resource.TestCheckResourceAttr(resourceName, "password_min_symbol", "1"),
					resource.TestCheckResourceAttr(resourceName, "password_exclude_username", "false"),
					resource.TestCheckResourceAttr(resourceName, "password_exclude_first_name", "true"),
					resource.TestCheckResourceAttr(resourceName, "password_exclude_last_name", "true"),
					resource.TestCheckResourceAttr(resourceName, "password_max_age_days", "60"),
					resource.TestCheckResourceAttr(resourceName, "password_expire_warn_days", "15"),
					resource.TestCheckResourceAttr(resourceName, "password_min_age_minutes", "60"),
					resource.TestCheckResourceAttr(resourceName, "password_history_count", "5"),
					resource.TestCheckResourceAttr(resourceName, "password_max_lockout_attempts", "3"),
					resource.TestCheckResourceAttr(resourceName, "password_auto_unlock_minutes", "2"),
					resource.TestCheckResourceAttr(resourceName, "password_show_lockout_failures", "true"),
					resource.TestCheckResourceAttr(resourceName, "question_min_length", "10"),
					resource.TestCheckResourceAttr(resourceName, "recovery_email_token", "20160"),
					resource.TestCheckResourceAttr(resourceName, "sms_recovery", "ACTIVE"),
				),
			},
		},
	})
}

func ensurePolicyExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		missingErr := fmt.Errorf("resource not found: %s", name)
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return missingErr
		}

		ID := rs.Primary.ID
		exist, err := doesPolicyExistsUpstream(ID)
		if err != nil {
			return err
		} else if !exist {
			return missingErr
		}

		return nil
	}
}

func createPolicyCheckDestroy(policyType string) func(*terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != policyType {
				continue
			}

			ID := rs.Primary.ID
			exists, err := doesPolicyExistsUpstream(ID)
			if err != nil {
				return err
			}

			if exists {
				return fmt.Errorf("policy still exists, ID: %s", ID)
			}
		}
		return nil
	}
}

func doesPolicyExistsUpstream(ID string) (bool, error) {
	client := getClientFromMetadata(testAccProvider.Meta())

	policy, resp, err := client.Policies.GetPolicy(ID)
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return policy.ID != "", nil
}
