package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccOktaAppFeatures_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppFeatures, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.example", resources.OktaIDaaSAppFeatures)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "USER_PROVISIONING"),
					resource.TestCheckResourceAttr(resourceName, "capabilities.create.lifecycle_create.status", "ENABLED"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "USER_PROVISIONING"),
					resource.TestCheckResourceAttr(resourceName, "capabilities.create.lifecycle_create.status", "DISABLED"),
				),
			},
		},
	})
}
