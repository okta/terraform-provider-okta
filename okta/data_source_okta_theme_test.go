package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOktaTheme_read(t *testing.T) {
	mgr := newFixtureManager(theme, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
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
