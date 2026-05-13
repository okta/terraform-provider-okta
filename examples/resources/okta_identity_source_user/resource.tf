resource "okta_identity_source_user" "test" {
  identity_source_id = "0oaxc95befZNgrJl71d7"
  id                 = "USEREXT123456TESTUSER2"

  profile {
    user_name  = "testuser2@example.com"
    email      = "testuser2@example.com"
    first_name = "Test"
    last_name  = "User2"
  }
}