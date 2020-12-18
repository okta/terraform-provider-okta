package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOktaTrustedOrigin_crud(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(trustedOrigin)
	config := mgr.GetFixtures("okta_trusted_origin.tf", ri, t)
	updatedConfig := mgr.GetFixtures("okta_trusted_origin_updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.testAcc_%d", trustedOrigin, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      testAccCheckTrustedOriginDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "origin", fmt.Sprintf("https://example2-%d.com", ri)),
					resource.TestCheckResourceAttr(resourceName, "active", "false"),
				),
			},
		},
	})
}

func testAccCheckTrustedOriginDestroy(s *terraform.State) error {
	client := getOktaClientFromMetadata(testAccProvider.Meta())

	for _, r := range s.RootModule().Resources {
		_, resp, err := client.TrustedOrigin.GetOrigin(context.Background(), r.Primary.ID)
		if is404(resp) {
			continue
		}
		if err != nil {
			return fmt.Errorf("error getting tructed origin: %v", err)
		}
		return fmt.Errorf("trusted origin still exists")
	}

	return nil
}
