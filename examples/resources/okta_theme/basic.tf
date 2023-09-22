# This example is part of the test harness. The okta_theme resource state has
# already been imported via import.tf

data "okta_brands" "test" {
}

resource "okta_theme" "example" {
  brand_id = tolist(data.okta_brands.test.brands)[0].id

  logo                                   = "../examples/okta_theme/okta_logo.png"
  favicon                                = "../examples/okta_theme/okta_favicon.png"
  background_image                       = "../examples/okta_theme/okta_background_image.png"
  primary_color_hex                      = "#1662dd"
  primary_color_contrast_hex             = "#ffffff"
  secondary_color_hex                    = "#ebebed"
  secondary_color_contrast_hex           = "#000000"
  sign_in_page_touch_point_variant       = "OKTA_DEFAULT"
  end_user_dashboard_touch_point_variant = "OKTA_DEFAULT"
  error_page_touch_point_variant         = "OKTA_DEFAULT"
  email_template_touch_point_variant     = "OKTA_DEFAULT"
}
