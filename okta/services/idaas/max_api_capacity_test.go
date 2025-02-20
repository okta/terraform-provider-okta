package idaas_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/provider"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccMaxApiCapacity_read(t *testing.T) {
	if provider.SkipVCRTest(t) {
		return
	}

	mgr := newFixtureManager("data-sources", resources.OktaIDaaSAppGroupAssignments, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	oldApiCapacity := os.Getenv("MAX_API_CAPACITY")
	t.Cleanup(func() {
		_ = os.Setenv("MAX_API_CAPACITY", oldApiCapacity)
	})
	// hack max api capacity value is enabled by env var
	os.Setenv("MAX_API_CAPACITY", "50")
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
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
