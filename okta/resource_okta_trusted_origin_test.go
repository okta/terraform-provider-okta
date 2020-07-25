package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOktaTrustedOrigin_crud(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager("okta_trusted_origin")
	config := mgr.GetFixtures("okta_trusted_origin.tf", ri, t)
	updatedConfig := mgr.GetFixtures("okta_trusted_origin_updated.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTrustedOriginDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						fmt.Sprintf("okta_trusted_origin.testAcc_%d", ri),
						"origin",
						fmt.Sprintf("https://example2-%d.com", ri),
					),
				),
			},
		},
	})
}

func testAccCheckTrustedOriginDestroy(s *terraform.State) error {
	ctx, client := getOktaClientFromMetadata(testAccProvider.Meta())

	for _, r := range s.RootModule().Resources {
		if _, resp, err := client.TrustedOrigin.GetOrigin(ctx, r.Primary.ID); err != nil {
			if resp.Status == "404 Not Found" {
				continue
			}
			return fmt.Errorf("[ERROR] Error Getting Trusted Origin in Okta: %v", err)
		}
		return fmt.Errorf("Trusted Origin still exists")
	}

	return nil
}
