---
page_title: "Data Source: okta_resource_owners"
description: |-
  Lists governance resources and their assigned owners.
---

# Data Source: okta_resource_owners

Lists governance resources and their assigned owners. Supports filtering by `parentResourceOrn` and `resource.orn`.

## Example Usage

```terraform
data "okta_resource_owners" "by_app" {
  filter = "parentResourceOrn eq \"orn:okta:idp:00o1234567890abcdef:apps:salesforce:0oa1234567890abcdef\""
}
```

## Argument Reference

- `filter` - (Required) A filter expression for listing resource owners. Supports `parentResourceOrn` (eq) and `resource.orn` (eq) filters.

## Attributes Reference

- `id` - Placeholder ID.
- `resource_owners` - List of resources with their assigned owners. Each element has:
    - `resource_id` - The ID of the resource.
    - `resource_type` - The type of the resource (e.g., `entitlement-bundles`).
    - `resource_orn` - The ORN of the resource.
    - `resource_name` - The name of the resource.
    - `parent_resource_orn` - The ORN of the parent resource (typically the app).
    - `principals` - List of principals (users or groups) that own the resource. Each element has:
        - `id` - The ID of the principal.
        - `type` - The type of the principal (`users` or `groups`).
        - `orn` - The ORN of the principal.
        - `name` - The display name of the principal.
