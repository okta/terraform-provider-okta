resource "okta_realm_assignment" "test" {
  name                 = "Example Realm Assignment"
  priority             = 55
  status               = "ACTIVE"
  profile_source_id    = okta_idp_saml.test.id
  condition_expression = "user.profile.login.contains(\"@example.com\")"
  realm_id             = okta_realm.example.id
}

resource "okta_realm" "example" {
  name       = "Example Realm"
  realm_type = "DEFAULT"
}
