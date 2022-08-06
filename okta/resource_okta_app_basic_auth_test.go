package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func TestAccAppBasicAuthApplication_crud(t *testing.T) {
	mgr := newFixtureManager(appBasicAuth, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", appBasicAuth)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appBasicAuth, createDoesAppExist(okta.NewBasicAuthApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewBasicAuthApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(mgr.Seed)),
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
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewBasicAuthApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(mgr.Seed)),
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
