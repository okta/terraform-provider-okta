---
layout: 'okta'
page_title: 'Okta: okta_behavior'
sidebar_current: 'docs-okta-datasource-behavior'
description: |- 
  Get a behavior by name or ID.
---

# okta_app

Use this data source to retrieve a behavior from Okta.

## Example Usage

```hcl
data "okta_behavior" "example" {
  label = "New City"
}
```

## Arguments Reference

- `name` - (Optional) The name of the behavior to retrieve. Name uses the `?q=<name>` query parameter exposed by 
  Okta's API.

- `id` - (Optional) `id` of behavior to retrieve, conflicts with `name`.

## Attributes Reference

- `id` - Behavior ID.

- `name` - Behavior name.

- `status` - Behavior status.

- `type` - Behavior type.

- `settings` - Map of behavior settings.
