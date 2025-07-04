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
  q = "testAcc_"
}

output "timestamp_fields_info" {
  value = {
    group_count                         = length(data.okta_groups.with_timestamps.groups)
    first_group_created                 = data.okta_groups.with_timestamps.groups[0].created
    first_group_last_updated            = data.okta_groups.with_timestamps.groups[0].last_updated
    first_group_last_membership_updated = data.okta_groups.with_timestamps.groups[0].last_membership_updated
  }
}
