---
page_title: "Resource: okta_resource_owner"
description: |-
  Manages an okta_resource_owner resource.
---

# Resource: okta_resource_owner

Manages an okta_resource_owner resource.

## Argument Reference

The following arguments are supported:
- `resource_orns` - (Required) The resources to assign owners
- `principal_orns` - (Optional) Owners of the resource.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

- `id` - The ID of the resource.
