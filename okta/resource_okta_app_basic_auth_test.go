package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/sdk"
)

func TestAccAppBasicAuthApplication_crud(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appBasicAuth)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appBasicAuth)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appBasicAuth, createDoesAppExist(sdk.NewBasicAuthApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewBasicAuthApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "url", "https://example.com/login.html"),
					resource.TestCheckResourceAttr(resourceName, "auth_url", "https://example.com/auth.html"),
					resource.TestCheckResourceAttr(resourceName, "groups.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewBasicAuthApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "url", "https://example.com/login.html"),
					resource.TestCheckResourceAttr(resourceName, "auth_url", "https://example.com/auth.html"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
				),
			},
		},
	})
}

func TestAccAppBasicAuthApplication_timeouts(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appBasicAuth)
	resourceName := fmt.Sprintf("%s.test", appBasicAuth)
	config := `
resource "okta_app_basic_auth" "test" {
  label    = "testAcc_replace_with_uuid"
  url      = "https://example.com/login.html"
  auth_url = "https://example.com/auth.html"
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
		CheckDestroy:      createCheckResourceDestroy(appBasicAuth, createDoesAppExist(sdk.NewBasicAuthApplication())),
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
