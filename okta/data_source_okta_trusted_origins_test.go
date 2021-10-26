package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDataSourceTrustedOrigin(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(trustedOrigins)
	config := mgr.GetFixtures("okta_trusted_origins.tf", ri, t)
	getTrustedOrigins := mgr.GetFixtures("datasource_okta_trusted_origins.tf", ri, t)
	resourceName := fmt.Sprintf("data.%s.test", trustedOrigins)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: getTrustedOrigins,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_auth_server.test", "id"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "claim_type", "IDENTITY"),
				),
			},
		},
	})
}
