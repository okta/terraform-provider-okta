package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaRateLimitWarningPercentage_cru(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSRateLimitWarningThresholdPercentage)
	mgr := newFixtureManager("resources", resources.OktaIDaaSRateLimitWarningThresholdPercentage, t.Name())
	config := mgr.GetFixtures("basic.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "warning_threshold", "90"),
				),
			},
		},
	})
}
