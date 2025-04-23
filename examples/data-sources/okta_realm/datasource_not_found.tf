# Should fail to find the realm doesn't exist
data "okta_realm" "test_not_found" {
  name = "Unknown Example Realm"
}
