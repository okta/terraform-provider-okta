# For testing TF_VAR_hostname is set in provider_test.go .
# In a live environment the operator would export `TF_VAR_hostname=[the
# hostname]` in order to expose hostname as a variable below.
variable "hostname" {
  type = string
}

resource "okta_resource_set" "test" {
  label       = "testAcc_replace_with_uuid"
  description = "testing, testing"
  resources = [
    "https://${var.hostname}/api/v1/users",
    "https://${var.hostname}/api/v1/apps",
    "https://${var.hostname}/api/v1/groups"
  ]
}

data "okta_org_metadata" "_" {}
locals {
  org_url = try(
    data.okta_org_metadata._.domains.alternate,
    data.okta_org_metadata._.domains.organization
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

### To Specify specific Groups
locals {
  org_url = "https://mycompany.okta.com"
}
resource "okta_resource_set" "test" {
  label       = "Specific Groups"
  description = "Only Specific Group"
  resources = [
    format("%s/api/v1/groups/groupid1", local.org_url),
    format("%s/api/v1/groups/groupid2", local.org_url)
  ]
}
