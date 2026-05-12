data "okta_identity_source_groups" "test" {
  identity_source_id = "0oaxc95befZNgrJl71d7"
  external_id        = "GROUPEXT123456784C2IF"
}

output "group_display_name" {
  value = data.okta_identity_source_groups.test.profile.display_name
}

output "group_description" {
  value = data.okta_identity_source_groups.test.profile.description
}