resource "okta_principal_rate_limit" "example" {
  principal_id = "<principal-id>"
  principal_type = "<principal_type>"

  # Optional fields
  # default_concurrency_percentage = 0
  # default_percentage = 0
}
