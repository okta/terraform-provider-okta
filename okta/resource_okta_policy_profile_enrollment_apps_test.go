package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaPolicyProfileEnrollmentApps(t *testing.T) {
	mgr := newFixtureManager(policyProfileEnrollmentApps, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", policyProfileEnrollmentApps)
	resourceName2 := fmt.Sprintf("%s.test_2", policyProfileEnrollmentApps)
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
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
