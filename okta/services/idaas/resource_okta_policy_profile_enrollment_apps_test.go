package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaPolicyProfileEnrollmentApps_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSPolicyProfileEnrollmentApps, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSPolicyProfileEnrollmentApps)
	resourceName2 := fmt.Sprintf("%s.test_2", resources.OktaIDaaSPolicyProfileEnrollmentApps)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
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
