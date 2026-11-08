resource "okta_realm_assignment" "test_optional" {
  name                 = "TestAcc Example Realm Assignment Optional"
  priority             = 50
  status               = "ACTIVE"
  condition_expression = "user.profile.login.contains(\"@optional.com\")"
  realm_id             = okta_realm.test_optional.id
}

resource "okta_realm" "test_optional" {
  name       = "TestAcc Example Assignment Realm Optional"
  realm_type = "DEFAULT"
}
