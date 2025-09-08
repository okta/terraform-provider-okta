# Should fail to find the realm assignment because it doesn't exist
data "okta_realm_assignment" "test_not_found" {
  name = "Unknown Example Realm Assignment"
}
