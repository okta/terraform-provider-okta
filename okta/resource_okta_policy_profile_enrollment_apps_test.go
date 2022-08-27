package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaPolicyProfileEnrollmentApps(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(policyProfileEnrollmentApps)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", policyProfileEnrollmentApps)
	resourceName2 := fmt.Sprintf("%s.test_2", policyProfileEnrollmentApps)
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					ensurePolicyExists(resourceName2),
					resource.TestCheckResourceAttr(resourceName, "apps.#", "1"),
					resource.TestCheckResourceAttr(resourceName2, "apps.#", "0"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					ensurePolicyExists(resourceName2),
					resource.TestCheckResourceAttr(resourceName, "apps.#", "0"),
					resource.TestCheckResourceAttr(resourceName2, "apps.#", "1"),
				),
			},
		},
	})
}
