resource "okta_realm" "test" {
  name       = "AccTest Example Realm"
  realm_type = "DEFAULT"
}

data "okta_realm" "test" {
  name       = "AccTest Example Realm"
  depends_on = [okta_realm.test]
}
