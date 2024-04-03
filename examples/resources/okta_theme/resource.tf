data "okta_brands" "test" {
}

# resource has been imported into current state:
# $ terraform import okta_theme.example <theme id>
resource "okta_theme" "example" {
  brand_id                               = tolist(data.okta_brands.test.brands)[0].id
  logo                                   = "path/to/logo.png"
  favicon                                = "path/to/favicon.png"
  background_image                       = "path/to/background.png"
  primary_color_hex                      = "#1662dd"
  secondary_color_hex                    = "#ebebed"
  sign_in_page_touch_point_variant       = "OKTA_DEFAULT"
  end_user_dashboard_touch_point_variant = "OKTA_DEFAULT"
  error_page_touch_point_variant         = "OKTA_DEFAULT"
  email_template_touch_point_variant     = "OKTA_DEFAULT"
}
