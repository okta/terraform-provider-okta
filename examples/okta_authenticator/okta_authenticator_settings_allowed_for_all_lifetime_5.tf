resource "okta_authenticator" "test" {
  type = "security_question"
  # TODO updating settings on the resource is not implemented yet
  settings = <<JSON
{
    "allowedFor": "all",
    "tokenLifetimeInMinutes": 5
}
JSON
}
