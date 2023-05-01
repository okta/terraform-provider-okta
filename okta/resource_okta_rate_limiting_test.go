package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaRateLimiting_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.example", rateLimiting)
	mgr := newFixtureManager(rateLimiting, t.Name())
	config := mgr.GetFixtures("basic.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
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
