---
layout: "okta"
page_title: "Okta: okta_groups"
sidebar_current: "docs-okta-datasource-groups"
description: |- Get a list of groups from Okta.
---

# okta_groups

Use this data source to retrieve a list of groups from Okta.

## Example Usage

```hcl
data "okta_groups" "example" {
  q = "Engineering - "
}
```

## Arguments Reference

- `q` - (Optional) Searches the name property of groups for matching value.

- `search` - (Optional) Searches for groups with a
  supported [filtering](https://developer.okta.com/docs/reference/api-overview/#filtering) expression for
  all [attributes](https://developer.okta.com/docs/reference/api/groups/#group-attributes)
  except for `"_embedded"`, `"_links"`, and `"objectClass"`

- `type` - (Optional) type of the group to retrieve. Can only be one of `OKTA_GROUP` (Native Okta Groups), `APP_GROUP`
  (Imported App Groups), or `BUILT_IN` (Okta System Groups).

## Attributes Reference

- `groups` - collection of groups retrieved from Okta with the following properties.
    - `id` - Group ID.
    - `name` - Group name.
    - `description` - Group description.
    - `type` - Group type.
