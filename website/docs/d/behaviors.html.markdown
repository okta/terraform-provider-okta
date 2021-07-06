---
layout: 'okta'
page_title: 'Okta: okta_behaviors'
sidebar_current: 'docs-okta-datasource-behaviors'
description: |- 
  Get a behaviors by search criteria.
---

# okta_app

Use this data source to retrieve a behaviors from Okta.

## Example Usage

```hcl
data "okta_behaviors" "example" {
  q = "New"
}
```

## Arguments Reference

- `q` - (Optional) Searches query to look up behaviors.

## Attributes Reference

- `behaviors` - List of behaviors.
  - `id` - Behavior ID.
  - `name` - Behavior name.
  - `status` - Behavior status.
  - `type` - Behavior type.
  - `settings` - Map of behavior settings.
