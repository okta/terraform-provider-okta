---
layout: 'okta'
page_title: 'Okta: okta_app'
sidebar_current: 'docs-okta-datasource-app'
description: |-
  Get an application of any kind from Okta.
---

# okta_app

Use this data source to retrieve the collaborators for a given repository.

## Example Usage

```hcl
data "okta_app" "example" {
  label = "Example App"
}
```

## Arguments Reference

- `label` - (Optional) The label of the app to retrieve, conflicts with `label_prefix` and `id`. Label uses the `?q=<label>` query parameter exposed by Okta's API. It should be noted that at this time this searches both `name` and `label`. This is used to avoid paginating through all applications.

- `label_prefix` - (Optional) Label prefix of the app to retrieve, conflicts with `label` and `id`. This will tell the provider to do a `starts with` query as opposed to an `equals` query.

- `id` - (Optional) `id` of application to retrieve, conflicts with `label` and `label_prefix`.

- `active_only` - (Optional) tells the provider to query for only `ACTIVE` applications.

## Attributes Reference

- `id` - `id` of application.

- `label` - `label` of application.

- `name` - `name` of application.

- `status` - `status` of application.
