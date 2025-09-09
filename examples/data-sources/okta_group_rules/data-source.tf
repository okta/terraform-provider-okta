data "okta_group_rules" "example" {
  search = "Engineering"
}

# Example with expand parameter to get group names
data "okta_group_rules" "example_with_expand" {
  search = "Engineering"
  expand = "groupIdToGroupNameMap"
}
