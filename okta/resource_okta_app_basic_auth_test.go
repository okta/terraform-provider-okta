package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/okta/okta-sdk-golang/okta"
)

func TestAccAppBasicAuthApplication_crud(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appBasicAuth)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appBasicAuth)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(appBasicAuth, createDoesAppExist(okta.NewBasicAuthApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewBasicAuthApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "url", "https://example.com/login.html"),
					resource.TestCheckResourceAttr(resourceName, "auth_url", "https://example.com/auth.html"),
					resource.TestCheckResourceAttr(resourceName, "groups.#", "1"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewBasicAuthApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "url", "https://example.com/login.html"),
					resource.TestCheckResourceAttr(resourceName, "auth_url", "https://example.com/auth.html"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "1"),
				),
			},
		},
	})
}
