locals {
  org_url = "https://mycompany.okta.com"
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
