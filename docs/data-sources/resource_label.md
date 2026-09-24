---
page_title: "Data Source: okta_resource_label"
description: |-
  Manages an okta_resource_label data source.
---

# Data Source: okta_resource_label

Lists all labeled resources  > **Note:** If you create a custom admin role to view labeled resources, ensure that the custom role has permissions to view the resource and governance labels.

## Argument Reference

The following arguments are supported:

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

- `id` - The ID of the data source.
- `labels` - List of assigned labels
- `orn` - The Okta resource, in [ORN format](https://developer.okta.com/docs/api/openapi/okta-management/guides/roles/#okta-resource-name-orn) format.
- `profile` - A limited set of properties from the resource's profile
