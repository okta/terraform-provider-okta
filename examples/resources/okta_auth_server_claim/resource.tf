resource "okta_auth_server_claim" "example" {
  auth_server_id = "<auth server id>"
  name           = "staff"
  value          = "String.substringAfter(user.email, \"@\") == \"example.com\""
  scopes         = ["${okta_auth_server_scope.example.name}"]
  claim_type     = "IDENTITY"
}
