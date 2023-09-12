resource "okta_group" "test_1" {
  name        = "testAcc_replace_with_uuid - Test 1"
  description = "testing, testing"
}

resource "okta_group" "test_2" {
  name        = "testAcc_replace_with_uuid  - Test 2"
  description = "testing, testing"
}

data "okta_groups" "test" {
  q = "testAcc_"
}

data "okta_groups" "app_groups" {
  type = "APP_GROUP"
}

output "special_groups" {
  # This is an example of syntax only. Okta API only adds group of type OKTA_GROUP
  # and so the example resource groups in this example above will be of that type.
  # OKTA_GROUP groups only have "name" and "description" properties which are
  # removed from okta-sdk-golang GroupProfileMap that is used to populate the bare
  # JSON custom_profile_attributes on the group data.

  # Groups of type APP_GROUP have customizable profile properties and a more
  # meaningful lookup example could be done like so:
  value = join(",", [for group in data.okta_groups.app_groups.groups : group.name
  if lookup(jsondecode(group.custom_profile_attributes), "some_attribute", "") == "Some Value"])
}
