package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaTheme_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSTheme, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_theme.test", "id"),
					resource.TestCheckResourceAttrSet("data.okta_theme.test", "logo_url"),
					resource.TestCheckResourceAttrSet("data.okta_theme.test", "favicon_url"),
					// resource.TestCheckResourceAttrSet("data.okta_theme.test", "background_image_url"), // background image is null on new orgs, skip check
					resource.TestCheckResourceAttrSet("data.okta_theme.test", "primary_color_hex"),
					resource.TestCheckResourceAttrSet("data.okta_theme.test", "primary_color_contrast_hex"),
					resource.TestCheckResourceAttrSet("data.okta_theme.test", "secondary_color_hex"),
					resource.TestCheckResourceAttrSet("data.okta_theme.test", "secondary_color_contrast_hex"),
					resource.TestCheckResourceAttrSet("data.okta_theme.test", "sign_in_page_touch_point_variant"),
					resource.TestCheckResourceAttrSet("data.okta_theme.test", "end_user_dashboard_touch_point_variant"),
					resource.TestCheckResourceAttrSet("data.okta_theme.test", "error_page_touch_point_variant"),
					resource.TestCheckResourceAttrSet("data.okta_theme.test", "email_template_touch_point_variant"),
					resource.TestCheckResourceAttrSet("data.okta_theme.test", "links"),
				),
			},
		},
	})
}
