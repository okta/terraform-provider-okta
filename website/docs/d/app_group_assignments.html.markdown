---
layout: 'okta'
page_title: 'Okta: okta_app_group_assignments'
sidebar_current: 'docs-okta-datasource-app-group-assignments'
description: |-
Get a set of groups assigned to an Okta application.
---


# okta_app_group_assignments

Use this data source to retrieve the list of groups assigned to the given Okta application (by ID).

## Example Usage

```hcl
data "okta_app_group_assignments" "test" {
  id = okta_app_oauth.test.id
}
```

## Argument Reference

- `id` - (Required) The ID of the Okta application you want to retrieve the groups for.

## Attribute Reference

- `id` - ID of application.

- `groups` - List of groups IDs assigned to the application.
