package okta

import (
	"fmt"
	"strings"
	"testing"

	"github.com/okta/okta-sdk-golang/okta"

	articulateOkta "github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func deletePasswordPolicies(artClient *articulateOkta.Client, client *okta.Client) error {
	return deletePolicyByType(passwordPolicyType, artClient, client)
}

func deletePolicyByType(t string, artClient *articulateOkta.Client, client *okta.Client) error {
	col, _, err := artClient.Policies.GetPoliciesByType(t)

	if err != nil {
		return fmt.Errorf("Failed to retrieve policies in order to properly destroy. Error: %s", err)
	}

	for _, policy := range col.Policies {
		if strings.HasPrefix(policy.Name, testResourcePrefix) {
			_, err = artClient.Policies.DeletePolicy(policy.ID)
		}
	}

	return nil
}

func TestAccOktaPolicyPassword(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicyPassword(ri)
	updatedConfig := testOktaPolicyPasswordUpdated(ri)
	resourceName := buildResourceFQN(passwordPolicy, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createPolicyCheckDestroy(passwordPolicy),
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

	policy, _, err := client.Policies.GetPolicy(ID)
	if is404(client) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return policy.ID != "", nil
}

func testOktaPolicyPassword(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  name        = "%s"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test Password Policy"
}
`, passwordPolicy, name, name)
}

// Noticed the below comment, added a TODO to actually validate for this.
// cannot change skipunlock to "true" if the authprovider is OKTA
// unless PASSWORD_POLICY_SOFT_LOCK is enabled
// (not supported in this TF provider at this time)
func testOktaPolicyPasswordUpdated(rInt int) string {
	name := buildResourceName(rInt)

	// Adding another resource so I can ensure the priority preference works
	return fmt.Sprintf(`
data "okta_everyone_group" "everyone-%d" {}

resource "%s" "%s" {
	name        = "%s"
	status      = "INACTIVE"
	description = "Terraform Acceptance Test Password Policy Updated"
	groups_included = [ "${data.okta_everyone_group.everyone-%d.id}" ]
	password_min_length = 12
	password_min_lowercase = 0 
	password_min_uppercase = 0 
	password_min_number = 0
	password_min_symbol = 1 
	password_exclude_username = false
	password_exclude_first_name = true 
	password_exclude_last_name = true 
	password_max_age_days = 60
	password_expire_warn_days = 15
	password_min_age_minutes = 60
	password_history_count = 5
	password_max_lockout_attempts = 3
	password_auto_unlock_minutes = 2
	password_show_lockout_failures = true
	question_min_length = 10
	recovery_email_token = 20160
	sms_recovery = "ACTIVE"
}
`, rInt, passwordPolicy, name, name, rInt)
}
