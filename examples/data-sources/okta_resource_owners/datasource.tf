data "okta_resource_owners" "by_app" {
  filter = "parentResourceOrn eq \"orn:okta:idp:00o1234567890abcdef:apps:salesforce:0oa1234567890abcdef\""
}

output "resource_owners" {
  value = data.okta_resource_owners.by_app.resource_owners
}
