resource "okta_application_feature_inbound_provisioning" "example" {
  app_id = "<app-id>"
  name = "Example Name"
  expression = "<expression>"
  expression = "<expression>"
  username_format = "Example Username Format"

  # Optional fields
  # allow_partial_match = true
  # auto_activate_new_users = true
  # auto_confirm_exact_match = true
  # auto_confirm_new_users = true
  # auto_confirm_partial_match = true
}
