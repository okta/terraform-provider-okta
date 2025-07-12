package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaFeatures_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSFeature, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSFeature)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "feature_id"),
					resource.TestCheckNoResourceAttr(resourceName, "mode"),
					resource.TestCheckNoResourceAttr(resourceName, "life_cycle"),
					resource.TestCheckResourceAttr(resourceName, "name", "Android Device Trust"),
					resource.TestCheckResourceAttr(resourceName, "description", "Collect a deeper set of device posture signals for Device Assurance by leveraging native Android Device Trust capabilities."),
					resource.TestCheckResourceAttr(resourceName, "status", "DISABLED"),
					resource.TestCheckResourceAttr(resourceName, "type", "self-service"),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "feature_id"),
					resource.TestCheckNoResourceAttr(resourceName, "mode"),
					resource.TestCheckResourceAttr(resourceName, "life_cycle", "ENABLE"),
					resource.TestCheckResourceAttr(resourceName, "name", "Android Device Trust"),
					resource.TestCheckResourceAttr(resourceName, "description", "Collect a deeper set of device posture signals for Device Assurance by leveraging native Android Device Trust capabilities."),
					resource.TestCheckResourceAttr(resourceName, "status", "ENABLED"),
					resource.TestCheckResourceAttr(resourceName, "type", "self-service"),
				),
			},
		},
	})
}
