# Example configuration showing how to use the okta_resource_sets data source
# This example demonstrates how to retrieve all resource sets in the organization

# Get all resource sets
data "okta_resource_sets" "all" {
}

# Example output showing the resource sets
output "resource_sets_count" {
  description = "The number of resource sets found"
  value       = length(data.okta_resource_sets.all.resource_sets)
}

output "resource_sets" {
  description = "List of all resource sets"
  value       = data.okta_resource_sets.all.resource_sets
}

# Example of filtering resource sets by label (if needed in the future)
# This would require adding a filter parameter to the datasource
/*
data "okta_resource_sets" "filtered" {
  label_prefix = "prod-"
}
*/
