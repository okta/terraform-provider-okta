---
page_title: "Data Source: okta_resource_owner"
description: |-
  Manages an okta_resource_owner data source.
---

# Data Source: okta_resource_owner

Lists all resources with assigned owners for an app (the parent resource).

## Argument Reference

The following arguments are supported:

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

- `id` - The ID of the data source.
- `parent_resource_orn` - The Okta resource, in [ORN format](https://developer.okta.com/docs/api/openapi/okta-management/guides/roles/#okta-resource-name-orn) format.
- `principals` - The principals that own the resource (users or groups)
- `resource` - Details of a resource that are owned by the principal, such as an app, an entitlement value, an entitlement bundle, or a collection
