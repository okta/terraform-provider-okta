package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaRoleSubscription_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", roleSubscription)
	mgr := newFixtureManager(roleSubscription, t.Name())
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
					resource.TestCheckResourceAttr(resourceName, "status", "unsubscribed")),
			},
		},
	})
}
