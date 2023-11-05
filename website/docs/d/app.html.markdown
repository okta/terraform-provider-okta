---
layout: 'okta' 
page_title: 'Okta: okta_app' 
sidebar_current: 'docs-okta-datasource-app' 
description: |- 
  Get an application of any kind from Okta.
---

# okta_app

Use this data source to retrieve an application from Okta.

## Example Usage

```hcl
data "okta_app" "example" {
  label = "Example App"
}
```

## Arguments Reference

- `label` - (Optional) The label or name of the app to retrieve, conflicts with `label_prefix` and `id`. Label uses
  the `?q=<label>` query parameter exposed by Okta's API. It should be noted that at this time the API searches both `name`
  and `label` with a [starts with query](https://developer.okta.com/docs/reference/api/apps/#list-applications) which
  may result in multiple apps being returned for the query. The data source further inspects the labels looking for
  an exact match.

- `label_prefix` - (Optional) The label or name prefix of the app to retrieve, conflicts with `label` and `id`.
  This will tell the provider to do a `starts with` query as opposed to an `equals` query.

- `id` - (Optional) `id` of application to retrieve, conflicts with `label` and `label_prefix`.

- `active_only` - (Optional) tells the provider to query for only `ACTIVE` applications.

## Attributes Reference

- `id` - Application ID.

- `label` - Application label.

- `name` - Application name.

- `status` - Application status.

- `links` - Generic JSON containing discoverable resources related to the app.
