package idaas_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaTheme_existing_update(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSTheme, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("updated.tf", t)
	importConfig := mgr.GetFixtures("import.tf", t)
	deleteImagesConfig := mgr.GetFixtures("delete-images.tf", t)

	// okta_theme is read and update only, so set up the test by importing the theme first
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				// this is set up only for import state test, ignore check as import.tf is for testing
				ExpectNonEmptyPlan: true,
				Config:             importConfig,
			},
			{
				ResourceName: "okta_theme.example",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["okta_theme.example"]
					if !ok {
						return "", fmt.Errorf("failed to find %s", "okta_theme.example")
					}

					brandID := rs.Primary.Attributes["brand_id"]
					themeID := rs.Primary.Attributes["theme_id"]

					return fmt.Sprintf("%s/%s", brandID, themeID), nil
				},
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					// import should only net one theme
					if len(s) != 1 {
						return errors.New("failed to import into resource into state")
					}
					// simple check
					if len(s[0].Attributes["links"]) <= 2 {
						return fmt.Errorf("there should more than two attributes set on the instance %+v", s[0].Attributes)
					}
					return nil
				},
			},
			{
				Config:  config,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_theme.example", "id"),
					resource.TestCheckResourceAttrSet("okta_theme.example", "logo_url"),
					resource.TestCheckResourceAttrSet("okta_theme.example", "favicon_url"),
					// resource.TestCheckResourceAttrSet("okta_theme.example", "background_image_url"), // background image is null on new orgs, skip check
					resource.TestCheckResourceAttrSet("okta_theme.example", "primary_color_hex"),
					resource.TestCheckResourceAttr("okta_theme.example", "primary_color_hex", "#1662dd"),
					resource.TestCheckResourceAttrSet("okta_theme.example", "primary_color_contrast_hex"),
					resource.TestCheckResourceAttrSet("okta_theme.example", "secondary_color_hex"),
					resource.TestCheckResourceAttr("okta_theme.example", "secondary_color_hex", "#ebebed"),
					resource.TestCheckResourceAttrSet("okta_theme.example", "secondary_color_contrast_hex"),
					resource.TestCheckResourceAttrSet("okta_theme.example", "sign_in_page_touch_point_variant"),
					resource.TestCheckResourceAttr("okta_theme.example", "sign_in_page_touch_point_variant", "OKTA_DEFAULT"),
					resource.TestCheckResourceAttrSet("okta_theme.example", "end_user_dashboard_touch_point_variant"),
					resource.TestCheckResourceAttr("okta_theme.example", "end_user_dashboard_touch_point_variant", "OKTA_DEFAULT"),
					resource.TestCheckResourceAttrSet("okta_theme.example", "error_page_touch_point_variant"),
					resource.TestCheckResourceAttr("okta_theme.example", "error_page_touch_point_variant", "OKTA_DEFAULT"),
					resource.TestCheckResourceAttrSet("okta_theme.example", "email_template_touch_point_variant"),
					resource.TestCheckResourceAttr("okta_theme.example", "email_template_touch_point_variant", "OKTA_DEFAULT"),
					resource.TestCheckResourceAttrSet("okta_theme.example", "links"),
				),
			},
			{
				Config:  updatedConfig,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_theme.example", "id"),
					resource.TestCheckResourceAttrSet("okta_theme.example", "logo_url"),
					resource.TestCheckResourceAttrSet("okta_theme.example", "favicon_url"),
					// resource.TestCheckResourceAttrSet("okta_theme.example", "background_image_url"), // background_image_url not present when dashboard is OKTA_DEFAULT
					resource.TestCheckResourceAttrSet("okta_theme.example", "primary_color_hex"),
					resource.TestCheckResourceAttr("okta_theme.example", "primary_color_hex", "#1662ff"),
					resource.TestCheckResourceAttrSet("okta_theme.example", "primary_color_contrast_hex"),
					resource.TestCheckResourceAttrSet("okta_theme.example", "secondary_color_hex"),
					resource.TestCheckResourceAttr("okta_theme.example", "secondary_color_hex", "#fbfbfd"),
					resource.TestCheckResourceAttrSet("okta_theme.example", "secondary_color_contrast_hex"),
					resource.TestCheckResourceAttrSet("okta_theme.example", "sign_in_page_touch_point_variant"),
					resource.TestCheckResourceAttr("okta_theme.example", "sign_in_page_touch_point_variant", "OKTA_DEFAULT"),
					resource.TestCheckResourceAttrSet("okta_theme.example", "end_user_dashboard_touch_point_variant"),
					resource.TestCheckResourceAttr("okta_theme.example", "end_user_dashboard_touch_point_variant", "OKTA_DEFAULT"),
					resource.TestCheckResourceAttrSet("okta_theme.example", "error_page_touch_point_variant"),
					resource.TestCheckResourceAttr("okta_theme.example", "error_page_touch_point_variant", "OKTA_DEFAULT"),
					resource.TestCheckResourceAttrSet("okta_theme.example", "email_template_touch_point_variant"),
					resource.TestCheckResourceAttr("okta_theme.example", "email_template_touch_point_variant", "OKTA_DEFAULT"),
					resource.TestCheckResourceAttrSet("okta_theme.example", "links"),
				),
			},
			{
				Config:  deleteImagesConfig,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_theme.example", "id"),

					// TODO find a better way to test delete of logo, favicon, background image, given:
					// Deletes a Theme logo. The org then uses the Okta default logo.
					// Deletes a Theme favicon. The org then uses the Okta default favicon.
					// Deletes a Theme background image
					// However, this still has fair coverage of delete as if their a failure the plan will fail.
					/*
						resource.TestCheckResourceAttrSet("okta_theme.example", "logo_url"),
						resource.TestCheckResourceAttr("okta_theme.example", "logo_url", ""),
						resource.TestCheckResourceAttrSet("okta_theme.example", "favicon_url"),
						resource.TestCheckResourceAttr("okta_theme.example", "favicon_url", ""),
						resource.TestCheckResourceAttrSet("okta_theme.example", "background_image_url"),
						resource.TestCheckResourceAttr("okta_theme.example", "background_image_url", ""),
					*/

					resource.TestCheckResourceAttrSet("okta_theme.example", "links"),
				),
			},
		},
	})
}
