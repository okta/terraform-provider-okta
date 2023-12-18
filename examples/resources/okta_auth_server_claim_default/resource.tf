resource "okta_auth_server_claim_default" "example" {
  auth_server_id = "<auth server id>"
  name           = "sub"
  value          = "(appuser != null) ? appuser.userName : app.clientId"
}
