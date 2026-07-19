resource "okta_user_client_token" "example" {
  user_id = "<user-id>"
  client_id = "<client-id>"

  # Optional fields
  # issuer = "<issuer>"
  # scopes = "<scopes>"
}
