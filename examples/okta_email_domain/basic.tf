# resource "okta_brand" "test" {
#   agree_to_custom_privacy_policy = true
#   custom_privacy_policy_url      = "https://example.com/privacy-policy"
#   remove_powered_by_okta         = false
# }

resource "okta_email_domain" "test" {
  brand_id     = "bnd5qwjvgpotf2LV51d7"
  domain       = "example.com"
  display_name = "test"
  user_name    = "fff"
}
