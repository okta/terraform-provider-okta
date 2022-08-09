package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOktaTrustedOrigin_crud(t *testing.T) {
	mgr := newFixtureManager(trustedOrigin, t.Name())
	config := mgr.GetFixtures("okta_trusted_origin.tf", t)
	updatedConfig := mgr.GetFixtures("okta_trusted_origin_updated.tf", t)
	resourceName := fmt.Sprintf("%s.testAcc_%d", trustedOrigin, mgr.Seed)

	oktaResourceTest(t, resource.TestCase{
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
					resource.TestCheckResourceAttr(resourceName, "origin", fmt.Sprintf("https://example2-%d.com", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "active", "false"),
				),
			},
		},
	})
}

func testAccCheckTrustedOriginDestroy(s *terraform.State) error {
	client := oktaClientForTest()

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
