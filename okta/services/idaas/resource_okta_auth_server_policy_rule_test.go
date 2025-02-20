package idaas_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
)

func TestAccResourceOktaAuthServerPolicyRule_create(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAuthServerPolicyRule)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAuthServerPolicyRule, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
		CheckDestroy:      checkResourceDestroy(resources.OktaIDaaSAuthServer, authServerExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "name", "test"),
					resource.TestCheckResourceAttr(resourceName, "system", "false"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "name", "test_updated"),
					resource.TestCheckResourceAttr(resourceName, "system", "false"),
				),
			},
		},
	})
}

// TestAccResourceOktaAuthServerPolicyRule_priority_concurrency_bug is a test to
// reproduce and then fix a bug in the Okta service where it couldn't, at the
// time, elegantly handle current API calls to either update rule create or rule
// priority.
func TestAccResourceOktaAuthServerPolicyRule_priority_concurrency_bug(t *testing.T) {
	numRules := 10
	testPolicyRules := make([]string, numRules)
	// Test setup makes each policy rule dependent on the one before it.
	for i := 0; i < numRules; i++ {
		dependsOn := i - 1
		testPolicyRules[i] = testPolicyRule(i, dependsOn)
	}

	config := fmt.Sprintf(`
data "okta_group" "all" {
  name = "Everyone"
}
resource "okta_auth_server" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "test"
  audiences   = ["whatever.rise.zone"]
}
resource "okta_auth_server_policy" "test" {
  name             = "test"
  description      = "test"
  priority         = 1
  client_whitelist = ["ALL_CLIENTS"]
  auth_server_id   = okta_auth_server.test.id
}
%s`, strings.Join(testPolicyRules, ""))
	mgr := newFixtureManager("resources", resources.OktaIDaaSAuthServerPolicyRule, t.Name())
	acctest.OktaResourceTest(
		t, resource.TestCase{
			PreCheck:          acctest.AccPreCheck(t),
			ErrorCheck:        testAccErrorChecks(t),
			ProviderFactories: acctest.AccProvidersFactoriesForTest(),
			CheckDestroy:      checkResourceDestroy(resources.OktaIDaaSAuthServer, authServerExists),
			Steps: []resource.TestStep{
				{
					Config: mgr.ConfigReplace(config),
					Check: resource.ComposeTestCheckFunc(
						// Just check if policy rule 09 exists. We only care
						// about reproducing then fixing the 500 bug. If we make
						// it to this check the test didn't fail and the bug has
						// been fixed.
						resource.TestCheckResourceAttr("okta_auth_server_policy_rule.test_09", "name", "Test Policy Rule 09"),
					),
				},
			},
		})
}

func testPolicyRule(num, dependsOn int) string {
	var dependsOnStr string
	if dependsOn >= 0 {
		dependsOnStr = fmt.Sprintf("depends_on = [okta_auth_server_policy_rule.test_%02d]", dependsOn)
	}

	return fmt.Sprintf(`
resource "okta_auth_server_policy_rule" "test_%02d" {
  auth_server_id       = okta_auth_server.test.id
  policy_id            = okta_auth_server_policy.test.id
  status               = "ACTIVE"
  name                 = "Test Policy Rule %02d"
  priority             = %d 
  group_whitelist      = [data.okta_group.all.id]
  grant_type_whitelist = ["implicit"]
  %s
}`,
		num, num, num+1, dependsOnStr)
}
