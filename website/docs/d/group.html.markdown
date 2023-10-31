---
layout: "okta"
page_title: "Okta: okta_group"
sidebar_current: "docs-okta-datasource-group"
description: |- Get a group from Okta.
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

-> `id` and `name` arguments are exclusive of each other

- `id` - (Required in lieu of `name`) ID of the group. Conflicts with `"name"` and `"type"`.

- `name` - (Required in lieu of `id`) name of group to retrieve. 

-> Okta API treats `name` as a starts with query. Therefore a name argument "My" will match any group starting with "My" such as "My Group" and "My Office"

- `type` - (Optional) type of the group to retrieve. Can only be one of `OKTA_GROUP` (Native Okta Groups), `APP_GROUP`
  (Imported App Groups), or `BUILT_IN` (Okta System Groups).

- `include_users` - (Optional) whether to retrieve all member ids.

- `delay_read_seconds` - (Optional) Force delay of the group read by N seconds. Useful when eventual consistency of group information needs to be allowed for; for instance, when group rules are known to have been applied.

## Attributes Reference

- `id` - ID of group.

- `name` - name of group.

- `type` - type of group.

- `description` - description of group.

- `users` - user ids that are members of this group, only included if `include_users` is set to `true`.
