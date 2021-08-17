package okta

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestMaxApiCapacity(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appGroupAssignments)
	config := mgr.GetFixtures("datasource.tf", ri, t)

	old := os.Getenv("MAX_API_CAPACITY")
	defer func() {
		_ = os.Setenv("MAX_API_CAPACITY", old)
	}()
	// hack max api capacity value is enabled by env var
	os.Setenv("MAX_API_CAPACITY", "50")
	resource.Test(t, resource.TestCase{
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
