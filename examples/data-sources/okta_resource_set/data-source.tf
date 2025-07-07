# Example configuration showing how to use the okta_resource_set data source
# This example demonstrates how to retrieve a resource set and use its resources
# in lifecycle preconditions to prevent users from being granted admin over themselves

# NOTE: In a future major release, this datasource will be split into separate datasources
# for better alignment with Terraform best practices. See the migration guide for details.

# Get a resource set by ID
data "okta_resource_set" "example" {
  id = "rs1234567890abcdef"
}

# Example output showing the resource set details
output "resource_set_label" {
  description = "The label of the resource set"
  value       = data.okta_resource_set.example.label
}

output "resource_set_description" {
  description = "The description of the resource set"
  value       = data.okta_resource_set.example.description
}

output "resource_set_resources" {
  description = "The resources included in the resource set"
  value       = data.okta_resource_set.example.resources
}

output "resource_set_resources_orn" {
  description = "The ORN resources included in the resource set"
  value       = data.okta_resource_set.example.resources_orn
}

# Example of using the datasource in a lifecycle precondition
# (This would be used in an actual resource, not as a standalone example)
/*
resource "okta_admin_role_custom_assignments" "example" {
  members         = ["user1", "user2"]
  custom_role_id  = "role123"
  resource_set_id = data.okta_resource_set.example.id

  lifecycle {
    precondition {
      condition = all([
        for member in ["user1", "user2"] : 
        !can(regex(".*${member}.*", join(",", data.okta_resource_set.example.resources)))
      ])
      error_message = "Members of a resource set should not be granted admin over themselves."
    }
  }
}
*/ 
