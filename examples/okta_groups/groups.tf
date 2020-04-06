resource "okta_group" "test_1" {
  name        = "My Group - 1 - testAcc_replace_with_uuid"
  description = "testing, testing"
}

resource "okta_group" "test_2" {
  name        = "My Group - 2 - testAcc_replace_with_uuid"
  description = "testing, testing"
}
