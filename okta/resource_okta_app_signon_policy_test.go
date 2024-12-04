package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaAppSignOnPolicy_crud(t *testing.T) {
	mgr := newFixtureManager("resources", appSignOnPolicy, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	renamedConfig := mgr.GetFixtures("basic_renamed.tf", t)
	resourceName := fmt.Sprintf("%v.test", appSignOnPolicy)

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
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceNameWithPrefix("testAcc_Test_App", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "description", "The updated app signon policy used by our test app."),
				),
			},
			{
				Config: renamedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceNameWithPrefix("testAcc_Test_App_Renamed", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "description", "The app signon policy used by our test app."),
				),
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
