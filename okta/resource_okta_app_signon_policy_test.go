package okta

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceOktaAppSignOnPolicy_crud(t *testing.T) {
	mgr := newFixtureManager("resources", appSignOnPolicy, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	renamedConfig := mgr.GetFixtures("basic_renamed.tf", t)
	resourceName := fmt.Sprintf("%v.test", appSignOnPolicy)

	importConfig := mgr.ConfigReplace(`
resource "okta_app_signon_policy_rule" "test" {
  depends_on = [okta_app_signon_policy.test]
}`)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:                 testAccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: testAccMergeProvidersFactories,
		CheckDestroy:             checkPolicyDestroy(appSignOnPolicy),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceNameWithPrefix("testAcc_Test_App", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "description", "The app signon policy used by our test app."),
					resource.TestCheckResourceAttrSet(resourceName, "default_rule_id"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceNameWithPrefix("testAcc_Test_App", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "description", "The updated app signon policy used by our test app."),
					resource.TestCheckResourceAttrSet(resourceName, "default_rule_id"),
				),
			},
			{
				Config: renamedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceNameWithPrefix("testAcc_Test_App_Renamed", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "description", "The app signon policy used by our test app."),
					resource.TestCheckResourceAttrSet(resourceName, "default_rule_id"),
				),
			},
			// check the default_rule_id was set by looking up the real policy rule by id and check its access is set to ALLOW
			{
				Config:       fmt.Sprintf("%s\n%s", renamedConfig, importConfig),
				ResourceName: "okta_app_signon_policy_rule.test",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["okta_app_signon_policy.test"]
					if !ok {
						return "", fmt.Errorf("failed to find policy")
					}

					policyId := rs.Primary.Attributes["id"]
					defaultRuleId := rs.Primary.Attributes["default_rule_id"]

					return fmt.Sprintf("%s/%s", policyId, defaultRuleId), nil
				},
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return errors.New("failed to import schema into state")
					}
					instance := states[0]
					access := instance.Attributes["access"]
					if access != "ALLOW" {
						return fmt.Errorf("expected okta_app_signon_policy.test access to be set to ALLOW, got %s", access)
					}

					return nil
				},
			},
		},
	})
}

// TestAccResourceOktaAppSignOnPolicy_default_policy_rule_can_be_set_to_deny
// tests to see if catch_all = false results in a default policy rule of DENY
func TestAccResourceOktaAppSignOnPolicy_default_policy_rule_can_be_set_to_deny(t *testing.T) {
	resourceName := fmt.Sprintf("%v.test", appSignOnPolicy)
	mgr := newFixtureManager("resources", appSignOnPolicy, t.Name())
	config := `
resource "okta_app_signon_policy" "test" {
  name        = "testAcc_Test_App_replace_with_uuid"
  description = "App Sign-On Policy with Default Rule DENY"
  catch_all   = false
}`
	importConfig := `
resource "okta_app_signon_policy" "test" {
  name        = "testAcc_Test_App_replace_with_uuid"
  description = "App Sign-On Policy with Default Rule DENY"
  catch_all   = false
}
resource "okta_app_signon_policy_rule" "test" {
  depends_on = [okta_app_signon_policy.test]
}`
	oktaResourceTest(t, resource.TestCase{
		PreCheck:                 testAccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: testAccMergeProvidersFactories,
		CheckDestroy:             checkPolicyDestroy(appSignOnPolicy),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr("okta_app_signon_policy.test", "name", buildResourceNameWithPrefix("testAcc_Test_App", mgr.Seed)),
					resource.TestCheckResourceAttr("okta_app_signon_policy.test", "description", "App Sign-On Policy with Default Rule DENY"),
					resource.TestCheckResourceAttr("okta_app_signon_policy.test", "catch_all", "false"),
					resource.TestCheckResourceAttrSet("okta_app_signon_policy.test", "default_rule_id"),
				),
			},
			{
				Config:       mgr.ConfigReplace(importConfig),
				ResourceName: "okta_app_signon_policy_rule.test",
				ImportState:  true,

				// fyi, import state id func equivelent to
				// terraform import okta_app_signon_policy_rule.test [policy id]/[policy rule id]
				// and is dirived directly off the attributes on okta_app_signon_policy.test
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["okta_app_signon_policy.test"]
					if !ok {
						return "", fmt.Errorf("failed to find policy")
					}

					policyId := rs.Primary.Attributes["id"]
					defaultRuleId := rs.Primary.Attributes["default_rule_id"]

					return fmt.Sprintf("%s/%s", policyId, defaultRuleId), nil
				},
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return errors.New("failed to import schema into state")
					}
					instance := states[0]
					access := instance.Attributes["access"]
					if access != "DENY" {
						return fmt.Errorf("expected okta_app_signon_policy.test access to be set to DENY, got %s", access)
					}

					return nil
				},
			},
		},
	})
}

func TestAccResourceOktaAppSignOnPolicy_destroy(t *testing.T) {
	mgr := newFixtureManager("resources", groupSchemaProperty, t.Name())
	oktaResourceTest(t, resource.TestCase{
		PreCheck:                 testAccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: testAccMergeProvidersFactories,
		CheckDestroy:             checkOktaGroupSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(`
# We create a sign on policy and two apps that have that sign on policy as their
# authentication policy.
# Check that apps have the policy as their authenication policy.
data "okta_policy" "test" {
  name = "Any two factors"
  type = "ACCESS_POLICY"
}
resource "okta_app_signon_policy" "test" {
  name        = "testAcc_Policy_replace_with_uuid"
  description = "Sign On Policy"
  depends_on = [
    data.okta_policy.test
  ]
}
resource "okta_app_oauth" "test1" {
  label = "testAcc_App_replace_with_uuid"
  type                      = "web"
  grant_types               = ["authorization_code"]
  redirect_uris             = ["http://localhost:3000"]
  post_logout_redirect_uris = ["http://localhost:3000"]
  response_types            = ["code"]
  authentication_policy     = okta_app_signon_policy.test.id
  depends_on = [
    data.okta_policy.test
  ]
}
resource "okta_app_oauth" "test2" {
  label = "testAcc_App_replace_with_uuid"
  type                      = "web"
  grant_types               = ["authorization_code"]
  redirect_uris             = ["http://localhost:3000"]
  post_logout_redirect_uris = ["http://localhost:3000"]
  response_types            = ["code"]
  authentication_policy     = okta_app_signon_policy.test.id
  depends_on = [
    data.okta_policy.test
  ]
}
data "okta_app_signon_policy" "test1" {
	app_id = okta_app_oauth.test1.id
}
data "okta_app_signon_policy" "test2" {
	app_id = okta_app_oauth.test2.id
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.okta_app_signon_policy.test1", "id", "data.okta_app_signon_policy.test2", "id"),
					resource.TestCheckResourceAttrPair("okta_app_signon_policy.test", "id", "data.okta_app_signon_policy.test1", "id"),
				),
			},
			{
				Config: mgr.ConfigReplace(`

# We destroy the sign on policy then check that both apps have been assigned to
# the default system policy as their authenication policy.
data "okta_policy" "test" {
  name = "Any two factors"
  type = "ACCESS_POLICY"
}
resource "okta_app_oauth" "test1" {
  label = "testAcc_App_replace_with_uuid"
  type                      = "web"
  grant_types               = ["authorization_code"]
  redirect_uris             = ["http://localhost:3000"]
  post_logout_redirect_uris = ["http://localhost:3000"]
  response_types            = ["code"]
  depends_on = [
    data.okta_policy.test
  ]
}
resource "okta_app_oauth" "test2" {
  label = "testAcc_App_replace_with_uuid"
  type                      = "web"
  grant_types               = ["authorization_code"]
  redirect_uris             = ["http://localhost:3000"]
  post_logout_redirect_uris = ["http://localhost:3000"]
  response_types            = ["code"]
  depends_on = [
    data.okta_policy.test
  ]
}
data "okta_app_signon_policy" "testA" {
	app_id = okta_app_oauth.test1.id
}
data "okta_app_signon_policy" "testB" {
	app_id = okta_app_oauth.test2.id
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.okta_app_signon_policy.testA", "id", "data.okta_app_signon_policy.testB", "id"),
				),
			},
			{
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.okta_app_signon_policy.testA", "id", "data.okta_app_signon_policy.testB", "id"),
					resource.TestCheckResourceAttrPair("data.okta_policy.test", "id", "data.okta_app_signon_policy.testA", "id"),
				),
			},
		},
	})
}
