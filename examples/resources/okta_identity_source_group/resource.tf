resource "okta_identity_source_group" "test" {
  identity_source_id = "0oaxc95befZNgrJl71d7"
  external_id        = "GRPEXT123456TESTGROUP1"

  profile {
    display_name = "Test Engineering Group"
    description  = "A test group for identity source integration"
  }
}
