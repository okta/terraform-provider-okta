---
layout: 'okta'
page_title: 'Okta: okta_policy'
sidebar_current: 'docs-okta-datasource-policy'
description: |-
  Get a policy from Okta.
---

# okta_policy

Use this data source to retrieve a policy from Okta.

## Example Usage

```hcl
data "okta_policy" "example" {
  name = "Password Policy Example"
  type = "PASSWORD"
}
```

## Arguments Reference

- `name` - (Required) name of policy to retrieve.

- `type` - (Required) type of policy to retrieve.

## Attributes Reference

- `id` - id of policy.

- `name` - name of policy.

- `type` - type of policy.
