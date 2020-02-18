package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOktaAuthServerScope_crud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.test", authServerScope)
	mgr := newFixtureManager(authServerScope)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(authServer, authServerExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "consent", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "name", "test:something"),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "consent", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "name", "test:something"),
					resource.TestCheckResourceAttr(resourceName, "description", "test_updated"),
				),
			},
		},
	})
}
