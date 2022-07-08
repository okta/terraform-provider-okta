package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaRateLimiting_crud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.example", rateLimiting)
	mgr := newFixtureManager(rateLimiting)
	config := mgr.GetFixtures("basic.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "login", "ENFORCE"),
					resource.TestCheckResourceAttr(resourceName, "authorize", "ENFORCE"),
					resource.TestCheckResourceAttr(resourceName, "communications_enabled", "true"),
				),
			},
		},
	})
}
