package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func TestAccAppBookmarkApplication_crud(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appBookmark)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appBookmark)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(appBookmark, createDoesAppExist(okta.NewBookmarkApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewBookmarkApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "url", "https://test.com"),
					resource.TestCheckResourceAttr(resourceName, "groups.#", "1"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewBookmarkApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "url", "https://test.com"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "1"),
				),
			},
		},
	})
}
