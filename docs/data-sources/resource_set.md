# okta_resource_set

Use this data source to retrieve a resource set by ID. This data source allows you to retrieve the details of a resource set including its resources, which can be used in lifecycle preconditions to prevent users from being granted admin over themselves.

## Example Usage

```hcl
data "okta_resource_set" "example" {
  id = "rs1234567890abcdef"
}

# Use the resource set in a lifecycle precondition
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
```

## Argument Reference

- `id` - (Required) The ID of the resource set to retrieve.

## Attributes Reference

- `id` - The ID of the resource set.
- `label` - Unique name given to the Resource Set.
- `description` - A description of the Resource Set.
- `resources` - The endpoints that reference the resources included in the Resource Set.
- `resources_orn` - The orn(Okta Resource Name) that reference the resources included in the Resource Set.

## Notes

- This data source is particularly useful when you need to validate that users are not being granted admin access over resources that include themselves.
- The `resources` and `resources_orn` attributes will be populated based on the type of resources in the resource set.
- Use this data source in combination with lifecycle preconditions to enforce security policies.

## Migration Notice

**Breaking Changes Planned for Future Major Release**

In a future major release, this datasource will be split into separate datasources to better align with Terraform best practices:

### Current Structure (to be deprecated):

```hcl
data "okta_resource_set" "example" {
  id = "rs1234567890abcdef"
}
# Returns: id, label, description, resources, resources_orn
```

### Future Structure:

```hcl
# Basic resource set metadata
data "okta_resource_set" "example" {
  id = "rs1234567890abcdef"
}
# Returns: id, label, description, created, last_updated

# Resource set resources with native IDs
data "okta_resource_set_resources" "example" {
  resource_set_id = data.okta_resource_set.example.id
}
# Returns: resource_set_id, resources (with native IDs instead of extracted links)
```

### Key Changes:

1. **Native IDs**: Resources will use native resource IDs instead of extracting IDs from `_links` attributes
2. **Single API Endpoint**: Each datasource will correspond to a single API endpoint
3. **Better Error Handling**: More reliable resource identification and error handling
4. **Improved Performance**: Reduced API calls and more efficient data processing

### Migration Path:

When the breaking change is released, users will need to:

1. Update their Terraform configurations to use the new datasource structure
2. Replace link-based resource references with native ID references
3. Update any lifecycle preconditions that depend on the current resource structure

We recommend planning for this migration in advance and testing with the new structure when it becomes available.
