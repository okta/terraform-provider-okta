package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaAuthServerClaimDefault(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.test", authServerClaimDefault)
	mgr := newFixtureManager(authServerClaimDefault)
	config := mgr.GetFixtures("basic.tf", ri, t)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(authServer, authServerExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "name", "address"),
					resource.TestCheckResourceAttr(resourceName, "value_type", "SYSTEM"),
					resource.TestCheckResourceAttr(resourceName, "claim_type", "IDENTITY"),
				),
			},
		},
	})
}
