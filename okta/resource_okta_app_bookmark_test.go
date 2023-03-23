package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/sdk"
)

func TestAccAppBookmarkApplication_crud(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appBookmark)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appBookmark)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appBookmark, createDoesAppExist(sdk.NewBookmarkApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewBookmarkApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "url", "https://test.com"),
					resource.TestCheckResourceAttr(resourceName, "groups.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewBookmarkApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "url", "https://test.com"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
				),
			},
		},
	})
}

func TestAccAppBookmarkApplication_timeouts(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appBookmark)
	resourceName := fmt.Sprintf("%s.test", appBookmark)
	config := `
resource "okta_app_bookmark" "test" {
  label = "testAcc_replace_with_uuid"
  url   = "https://test.com"
  timeouts {
    create = "60m"
    read = "2h"
    update = "30m"
  }
}`
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appBookmark, createDoesAppExist(sdk.NewBookmarkApplication())),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config, ri),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewAutoLoginApplication())),
					resource.TestCheckResourceAttr(resourceName, "timeouts.create", "60m"),
					resource.TestCheckResourceAttr(resourceName, "timeouts.read", "2h"),
					resource.TestCheckResourceAttr(resourceName, "timeouts.update", "30m"),
				),
			},
		},
	})
}

// TestAccAppBookmarkApplication_PR1366 Test for @jakezarobsky-8451 PR #1366
// https://github.com/okta/terraform-provider-okta/pull/1366
func TestAccAppBookmarkApplication_PR1366_authentication_policy(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appBookmark)
	resourceName := fmt.Sprintf("%s.test", appBookmark)
	config := `
resource "okta_group" "group" {
  name = "testAcc_replace_with_uuid"
}
data "okta_policy" "test" {
  name = "Any two factors"
  type = "ACCESS_POLICY"
}
resource "okta_app_signon_policy" "test" {
  name        = "testAcc_Policy_replace_with_uuid"
  description = "Sign On Policy"
  depends_on  = [
    data.okta_policy.test
  ]
}
resource "okta_user" "user" {
  admin_roles = ["APP_ADMIN", "USER_ADMIN"]
  first_name  = "TestAcc"
  last_name   = "blah"
  login       = "testAcc-replace_with_uuid@example.com"
  email       = "testAcc-replace_with_uuid@example.com"
}
resource "okta_app_bookmark" "test" {
  label  = "testAcc_replace_with_uuid"
  url    = "https://test.com"
  groups = [okta_group.group.id]
  users {
    id       = okta_user.user.id
    username = okta_user.user.email
  }
  authentication_policy = okta_app_signon_policy.test.id
}`
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appBookmark, createDoesAppExist(sdk.NewBookmarkApplication())),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config, ri),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewAutoLoginApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "url", "https://test.com"),
					resource.TestCheckResourceAttr(resourceName, "groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
					resource.TestCheckResourceAttrSet(resourceName, "authentication_policy"),
				),
			},
		},
	})
}
