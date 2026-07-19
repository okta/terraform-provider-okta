resource "okta_user_factor_token_hotp" "example" {
  user_id = "<user-id>"
  factor_type = "<factor_type>"

  # Optional fields
  # factor_profile_id = "<factor-profile-id>"
  # profile = "<profile>"
  # provider = "<provider>"
}
