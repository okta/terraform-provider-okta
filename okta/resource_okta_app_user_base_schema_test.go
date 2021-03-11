package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func TestAccAppUserBaseSchema_change(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appUserBaseSchema)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updated := mgr.GetFixtures("updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appUserBaseSchema)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		// Just need to make sure the app gets cleaned up
		CheckDestroy: createCheckResourceDestroy(appOAuth, createDoesAppExist(okta.NewOpenIdConnectApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "index", "name"),
					resource.TestCheckResourceAttr(resourceName, "title", "Name"),
					resource.TestCheckResourceAttr(resourceName, "type", "string"),
					resource.TestCheckResourceAttr(resourceName, "required", "true"),
					resource.TestCheckResourceAttr(resourceName, "permissions", "READ_ONLY"),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "index", "name"),
					resource.TestCheckResourceAttr(resourceName, "title", "Name"),
					resource.TestCheckResourceAttr(resourceName, "type", "string"),
					resource.TestCheckResourceAttr(resourceName, "required", "false"),
					resource.TestCheckResourceAttr(resourceName, "permissions", "READ_ONLY"),
				),
			},
		},
	})
}
