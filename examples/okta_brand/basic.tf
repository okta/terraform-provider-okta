# This example is part of the test harness. The okta_brand resource state has
# already been imported via import.tf

resource "okta_brand" "example" {
  agree_to_custom_privacy_policy = true
  custom_privacy_policy_url      = "https://example.com/privacy-policy"
  remove_powered_by_okta         = false
}
