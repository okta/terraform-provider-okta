resource "okta_group" "test_1" {
  name        = "testAcc_replace_with_uuid - Test 1"
  description = "testing, testing"
}

resource "okta_group" "test_2" {
  name        = "testAcc_replace_with_uuid  - Test 2"
  description = "testing, testing"
}

# Test sorting by created date in ascending order
data "okta_groups" "sorted_by_created" {
  q          = "testAcc_"
  sort_by    = "created"
  sort_order = "asc"
}

# Test sorting by name in descending order
data "okta_groups" "sorted_by_name" {
  q          = "testAcc_"
  sort_by    = "name"
  sort_order = "desc"
}

# Test sorting by lastUpdated in descending order
data "okta_groups" "sorted_by_last_updated" {
  q          = "testAcc_"
  sort_by    = "lastUpdated"
  sort_order = "desc"
}

# Test sorting by ID in ascending order
data "okta_groups" "sort_by_id_asc" {
  q          = "testAcc_"
  sort_by    = "id"
  sort_order = "asc"
}

# Test sorting by ID in descending order
data "okta_groups" "sort_by_id_desc" {
  q          = "testAcc_"
  sort_by    = "id"
  sort_order = "desc"
}

# Test sorting by lastMembershipUpdated in ascending order
data "okta_groups" "sort_by_membership_updated" {
  q          = "testAcc_"
  sort_by    = "lastMembershipUpdated"
  sort_order = "asc"
}

# Test sorting with type filter combination
data "okta_groups" "sort_with_type_filter" {
  q          = "testAcc_"
  type       = "OKTA_GROUP"
  sort_by    = "created"
  sort_order = "desc"
}

output "comprehensive_sorting_info" {
  value = {
    # Basic sorting tests
    created_sorted_count = length(data.okta_groups.sorted_by_created.groups)
    name_sorted_count    = length(data.okta_groups.sorted_by_name.groups)
    last_updated_count   = length(data.okta_groups.sorted_by_last_updated.groups)

    # Additional sort field tests
    id_asc_count             = length(data.okta_groups.sort_by_id_asc.groups)
    id_desc_count            = length(data.okta_groups.sort_by_id_desc.groups)
    membership_updated_count = length(data.okta_groups.sort_by_membership_updated.groups)

    # Combination tests
    type_filter_count = length(data.okta_groups.sort_with_type_filter.groups)
  }
}
