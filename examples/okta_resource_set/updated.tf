resource "okta_resource_set" "test" {
  label       = "testAcc_replace_with_uuid"
  description = "testing, testing updated"
  resources   = [
    "https://terraform-provider-okta.oktapreview.com/api/v1/users",
    "https://terraform-provider-okta.oktapreview.com/api/v1/apps",
  ]
}
