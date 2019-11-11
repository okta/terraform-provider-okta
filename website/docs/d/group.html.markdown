---
layout: "okta"
page_title: "Okta: okta_group"
sidebar_current: "docs-okta-datasource-group"
description: |-
  Get a group from Okta.
---

# okta_group

Use this data source to retrieve a group from Okta.

## Example Usage

```hcl
data "okta_group" "example" {
  name = "Example App"
}
```

## Arguments Reference

* `name` - (Required) name of group to retrieve.

* `include_users` - (Optional) whether or not to retrieve all member ids.

## Attributes Reference

* `id` - id of group.

* `name` - name of group.

* `description` - description of group.

* `users` - user ids that are members of this group, only included if `include_users` is set to `true`.
