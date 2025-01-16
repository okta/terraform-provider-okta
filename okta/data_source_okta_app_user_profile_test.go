package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceAppUserProfile_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager("okta_app_user_profile")
	config := mgr.GetFixtures("datasource.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_app_user_profile.test", "profile"),
				),
			},
		},
	})
}
