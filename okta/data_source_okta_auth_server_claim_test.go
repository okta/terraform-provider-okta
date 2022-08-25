package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceAuthServerClaim(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(authServerClaim)
	config := mgr.GetFixtures("datasource.tf", ri, t)
	createUser := mgr.GetFixtures("datasource_create_auth_server.tf", ri, t)
	resourceName := fmt.Sprintf("data.%s.test", authServerClaim)
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: createUser,
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
