locals {
  org_url = "https://${var.hostname}"
}

resource "okta_resource_set" "test" {
  label       = "UsersAppsAndGroups"
  description = "All the users, app and groups"
  resources = [
    format("%s/api/v1/users", local.org_url),
    format("%s/api/v1/apps", local.org_url),
    format("%s/api/v1/groups", local.org_url)
  ]
}

data "okta_org_metadata" "_" {}
locals {
  org_url = try(
    data.okta_org_metadata._.alternate,
    data.okta_org_metadata._.organization
  )
}
resource "okta_resource_set" "example" {
  label       = "UsersAppsAndGroups"
  description = "All the users, app and groups"
  resources = [
    "${local.org_url}/api/v1/users",
    "${local.org_url}/api/v1/apps",
    "${local.org_url}/api/v1/groups"
  ]
}

### To Provide permissions to specific Groups

locals {
  org_url = "https://${var.hostname}"
}
resource "okta_resource_set" "test" {
  label       = "Specific Groups"
  description = "Only Specific Group"
  resources = [
    format("%s/api/v1/groups/groupid1", local.org_url),
    format("%s/api/v1/groups/groupid2", local.org_url)
  ]
}
