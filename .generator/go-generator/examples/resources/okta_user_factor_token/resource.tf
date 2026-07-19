resource "okta_user_factor_token" "example" {
  user_id = "<user-id>"
  factor_type = "<factor_type>"

  # Optional fields
  # profile = "<profile>"
  # provider = "<provider>"
  # verify = "<verify>"
}
