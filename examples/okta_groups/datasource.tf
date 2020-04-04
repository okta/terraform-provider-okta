resource "okta_group" "test_1" {
  name        = "My Group - testAcc_replace_with_uuid"
  description = "testing, testing"
}

resource "okta_group" "test_2" {
  name        = "My Group - testAcc_replace_with_uuid"
  description = "testing, testing"
}

data "okta_groups" "test" {
  q = "My Group - "

  depends_on = [
    okta_group.test_1,
    okta_group.test_2
  ]
}
