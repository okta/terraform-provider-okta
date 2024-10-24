---
layout: 'okta' 
page_title: 'Okta: okta_apps' 
sidebar_current: 'docs-okta-datasource-apps' 
description: |- 
  List applications of any kind from Okta.
---

# okta_app

Use this data source to list applications from Okta.

## Example Usage

```hcl
data "okta_apps" "example" {
  label_prefix = "Example App"
}
```

## Arguments Reference

- `label` - (Optional) The label or name of the apps to retrieve, conflicts with `label_prefix`. Both `label` and `label_prefix`
  use the `?q=<label>` query parameter exposed by Okta's API. It should be noted that at this time the API searches both `name`
  and `label` with a [starts with query](https://developer.okta.com/docs/reference/api/apps/#list-applications).
  The data source further inspects the labels looking for an exact match when the `label` parameter is set.

- `label_prefix` - (Optional) The label or name prefix of the apps to retrieve, conflicts with `label`. It should be noted that
  Okta's API searches both the `label` and `name` application properties.

- `active_only` - (Optional) tells the provider to query for only `ACTIVE` applications.

## Attributes Reference

- `apps` - collection of applications retrieved from Okta with the following properties.

  - `id` - Application ID.

  - `label` - Application label.

  - `name` - Application name.

  - `status` - Application status.

  - `links` - Generic JSON containing discoverable resources related to the app.
