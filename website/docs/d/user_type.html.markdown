---
layout: 'okta'
page_title: 'Okta: okta_user_type'
sidebar_current: 'docs-okta-datasource-user-type'
description: |-
  Get a user type from Okta.
---

# okta_group

Use this data source to retrieve a user type from Okta.

## Example Usage

```hcl
data "okta_user_type" "example" {
  name = "example"
}
```

## Arguments Reference

- `name` - (Required) name of user type to retrieve.

## Attributes Reference

- `id` - id of user type.

- `name` - name of user type.

- `display_name` - display name of user type.

- `description` - description of user type.
