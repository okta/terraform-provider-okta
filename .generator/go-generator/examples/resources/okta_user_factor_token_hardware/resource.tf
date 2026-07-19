resource "okta_user_factor_token_hardware" "example" {
  user_id = "<user-id>"
  factor_type = "<factor_type>"

  # Optional fields
  # profile = "<profile>"
  # provider = "<provider>"
  # pass_code = "<pass_code>"
}
