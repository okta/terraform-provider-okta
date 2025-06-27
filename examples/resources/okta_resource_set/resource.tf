# Example 1: Basic resource set with direct URLs
resource "okta_resource_set" "example" {
  label       = "Basic Resource Set"
  description = "Resource set for managing API endpoints"
  resources = [
    "https://your-domain.okta.com/api/v1/users",
    "https://your-domain.okta.com/api/v1/groups"
  ]
}

# Example 2: Using org metadata to get the correct domain
data "okta_org_metadata" "org" {}

resource "okta_resource_set" "dynamic" {
  label       = "Dynamic Resource Set"
  description = "Resource set using org metadata"
  resources = [
    "${data.okta_org_metadata.org.organization}/api/v1/users",
    "${data.okta_org_metadata.org.organization}/api/v1/groups"
  ]
}

# Example 3: Specific group access using ORN format
resource "okta_resource_set" "specific_groups" {
  label       = "Specific Groups"
  description = "Access to specific groups only"
  resources_orn = [
    "orn:okta:directory:${data.okta_org_metadata.org.id}:group:groupid1",
    "orn:okta:directory:${data.okta_org_metadata.org.id}:group:groupid2"
  ]
}
