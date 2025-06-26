package idaas_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func TestAccResourceOktaTrustedOrigin_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSTrustedOrigin, t.Name())
	config := mgr.GetFixtures("okta_trusted_origin.tf", t)
	updatedConfig := mgr.GetFixtures("okta_trusted_origin_updated.tf", t)
	resourceName := fmt.Sprintf("%s.testAcc_%d", resources.OktaIDaaSTrustedOrigin, mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkTrustedOriginDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "origin", fmt.Sprintf("https://example2-%d.com", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "active", "false"),
				),
			},
		},
	})
}

func checkTrustedOriginDestroy(s *terraform.State) error {
	client := iDaaSAPIClientForTestUtil.OktaSDKClientV2()

	for _, r := range s.RootModule().Resources {
		_, resp, err := client.TrustedOrigin.GetOrigin(context.Background(), r.Primary.ID)
		if utils.Is404(resp) {
			continue
		}
		if err != nil {
			return fmt.Errorf("error getting tructed origin: %v", err)
		}
		return fmt.Errorf("trusted origin still exists")
	}

	return nil
}
