---
page_title: "Data Source: okta_resource_owners_catalog_resource"
description: |-
  Manages an okta_resource_owners_catalog_resource data source.
---

# Data Source: okta_resource_owners_catalog_resource

Lists all resources without assigned owners for an app (the parent resource).

## Argument Reference

The following arguments are supported:

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

- `id` - The ID of the data source.
- `data` - Resource owner details
- `parent_resource_orn` - The Okta resource, in [ORN format](https://developer.okta.com/docs/api/openapi/okta-management/guides/roles/#okta-resource-name-orn) format.
