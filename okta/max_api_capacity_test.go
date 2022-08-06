package okta

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMaxApiCapacity(t *testing.T) {
	mgr := newFixtureManager(appGroupAssignments, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	oldApiCapacity := os.Getenv("MAX_API_CAPACITY")
	t.Cleanup(func() {
		_ = os.Setenv("MAX_API_CAPACITY", oldApiCapacity)
	})
	// hack max api capacity value is enabled by env var
	os.Setenv("MAX_API_CAPACITY", "50")
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_app_group_assignments.test", "groups.#"),
				),
			},
		},
	})
}
