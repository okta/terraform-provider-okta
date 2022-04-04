package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaAuthServerDefault_crud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.sun_also_rises", authServerDefault)
	mgr := newFixtureManager(authServerDefault)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, authServerExists),
					resource.TestCheckResourceAttr(resourceName, "description", "Default Authorization Server for your Applications"),
					resource.TestCheckResourceAttr(resourceName, "name", "default"),
					resource.TestCheckResourceAttr(resourceName, "audiences.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "credentials_rotation_mode", "MANUAL"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, authServerExists),
					resource.TestCheckResourceAttr(resourceName, "description", "Default Authorization Server"),
					resource.TestCheckResourceAttr(resourceName, "name", "default"),
					resource.TestCheckResourceAttr(resourceName, "audiences.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "credentials_rotation_mode", "MANUAL"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, authServerExists),
					resource.TestCheckResourceAttr(resourceName, "description", "Default Authorization Server for your Applications"),
					resource.TestCheckResourceAttr(resourceName, "name", "default"),
					resource.TestCheckResourceAttr(resourceName, "audiences.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "credentials_rotation_mode", "MANUAL"),
				),
			},
		},
	})
}
