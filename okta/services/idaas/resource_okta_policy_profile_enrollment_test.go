package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
)

func TestAccResourceOktaPolicyProfileEnrollment_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSPolicyProfileEnrollment, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSPolicyProfileEnrollment)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:      checkPolicyDestroy(resources.OktaIDaaSPolicyProfileEnrollment),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusInactive),
				),
			},
		},
	})
}
