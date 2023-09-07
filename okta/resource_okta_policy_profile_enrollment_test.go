package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaPolicyProfileEnrollment(t *testing.T) {
	mgr := newFixtureManager(policyProfileEnrollment, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", policyProfileEnrollment)
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkPolicyDestroy(policyProfileEnrollment),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusInactive),
				),
			},
		},
	})
}
