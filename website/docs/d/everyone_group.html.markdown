---
layout: 'okta'
page_title: 'Okta: okta_everyone_group'
sidebar_current: 'docs-okta-datasource-everyone-group'
description: |-
  Get the `Everyone` group from Okta.
---

# okta_everyone_group

Use this data source to retrieve the `Everyone` group from Okta. The same can be achieved with the `okta_group` data
source with `name = "Everyone"`. This is simply a shortcut.

## Example Usage

```hcl
data "okta_everyone_group" "example" {}
```

## Arguments Reference

- `include_users` - (Optional) whether to retrieve all member ids.

## Attributes Reference

- `id` - ID of the group.

- `description` - description of group.
