resource "okta_brand" "example" {
  agree_to_custom_privacy_policy = true
  custom_privacy_policy_url      = "https://example.com/privacy-policy"
  remove_powered_by_okta         = false
}
