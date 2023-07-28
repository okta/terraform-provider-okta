package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/sdk"
)

func TestAccAppBookmarkApplication_crud(t *testing.T) {
	mgr := newFixtureManager(appBookmark, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", appBookmark)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(appBookmark, createDoesAppExist(sdk.NewBookmarkApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewBookmarkApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "url", "https://test.com"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewBookmarkApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "url", "https://example.com"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
				),
			},
		},
	})
}

func TestAccAppBookmarkApplication_timeouts(t *testing.T) {
	mgr := newFixtureManager(appBookmark, t.Name())
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
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(appBookmark, createDoesAppExist(sdk.NewBookmarkApplication())),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
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
	mgr := newFixtureManager(appBookmark, t.Name())
	resourceName := fmt.Sprintf("%s.test", appBookmark)
	config := `
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
resource "okta_app_bookmark" "test" {
  label  = "testAcc_replace_with_uuid"
  url    = "https://test.com"
  authentication_policy = okta_app_signon_policy.test.id
}`
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(appBookmark, createDoesAppExist(sdk.NewBookmarkApplication())),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewAutoLoginApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "url", "https://test.com"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
					resource.TestCheckResourceAttrSet(resourceName, "authentication_policy"),
				),
			},
		},
	})
}
