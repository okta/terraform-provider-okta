resource "okta_group" "test_1" {
  name        = "testAcc_replace_with_uuid - Test 1"
  description = "testing, testing"
}

resource "okta_group" "test_2" {
  name        = "testAcc_replace_with_uuid  - Test 2"
  description = "testing, testing"
}

# Test timestamp fields exposure
data "okta_groups" "with_timestamps" {
  q = "testAcc_replace_with_uuid"
}


