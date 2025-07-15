# Example configuration showing how to use the okta_resource_set_resources data source
# This example demonstrates how to retrieve resources within a specific resource set

# Get resources within a specific resource set
data "okta_resource_set_resources" "example" {
  resource_set_id = "rs1234567890abcdef"
}

# Example output showing the resources
output "resources_count" {
  description = "The number of resources in the resource set"
  value       = length(data.okta_resource_set_resources.example.resources)
}

output "resources" {
  description = "List of resources in the resource set"
  value       = data.okta_resource_set_resources.example.resources
}

# Example of using the datasource with a resource set created by the provider
/*
resource "okta_resource_set" "example" {
  label       = "example-resource-set"
  description = "An example resource set"
  resources   = ["https://example.okta.com/api/v1/users", "https://example.okta.com/api/v1/groups"]
}

data "okta_resource_set_resources" "example" {
  resource_set_id = okta_resource_set.example.id
}
*/
