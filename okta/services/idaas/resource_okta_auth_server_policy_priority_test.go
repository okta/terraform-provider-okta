package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaAuthServerPolicyPriority_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAuthServerPolicyPriority)
	firstPolicy := fmt.Sprintf("%s.first", resources.OktaIDaaSAuthServerPolicy)
	secondPolicy := fmt.Sprintf("%s.second", resources.OktaIDaaSAuthServerPolicy)
	thirdPolicy := fmt.Sprintf("%s.third", resources.OktaIDaaSAuthServerPolicy)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAuthServerPolicyPriority, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAuthServer, authServerExists),
		Steps: []resource.TestStep{
			{
				// Step 1: create with order [first, second, third]
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "priorities.#", "3"),
					resource.TestCheckResourceAttrPair(resourceName, "priorities.0", firstPolicy, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "priorities.1", secondPolicy, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "priorities.2", thirdPolicy, "id"),
				),
			},
			{
				// Step 2: reorder to [third, first, second]
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "priorities.#", "3"),
					resource.TestCheckResourceAttrPair(resourceName, "priorities.0", thirdPolicy, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "priorities.1", firstPolicy, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "priorities.2", secondPolicy, "id"),
				),
			},
			{
				// Step 3: import — ignore priorities since import captures all policies on the
				// auth server, which may differ in order from what config expects post-import.
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"priorities"},
			},
		},
	})
}
