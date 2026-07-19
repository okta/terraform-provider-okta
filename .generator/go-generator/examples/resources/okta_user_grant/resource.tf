resource "okta_user_grant" "example" {
  user_id = "<user-id>"
  issuer = "<issuer>"
  scope_id = "<scope-id>"
}
