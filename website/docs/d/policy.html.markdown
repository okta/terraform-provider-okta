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

- `name` - (Required) Name of policy to retrieve.

- `type` - (Required) Type of policy to retrieve. Valid values: `OKTA_SIGN_ON`, `PASSWORD`, `MFA_ENROLL`, `IDP_DISCOVERY`

## Attributes Reference

- `id` - id of policy.

- `name` - name of policy.

- `type` - type of policy.
