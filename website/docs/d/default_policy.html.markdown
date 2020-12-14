---
layout: 'okta'
page_title: 'Okta: okta_default_policy'
sidebar_current: 'docs-okta-datasource-default-policy'
description: |-
  Get a Default policy from Okta.
---

# okta_default_policy

Use this data source to retrieve a "Default" policy from Okta. This same thing can be achieved using the `okta_policy` with `name = "Default"`, this is simply a shortcut.

## Example Usage

```hcl
data "okta_default_policy" "example" {
  type = "PASSWORD"
}
```

## Arguments Reference

- `type` - (Required) type of policy to retrieve.  Valid values: `OKTA_SIGN_ON`, `PASSWORD`, `MFA_ENROLL`, `OAUTH_AUTHORIZATION_POLICY`, `IDP_DISCOVERY`

## Attributes Reference

- `id` - id of policy.

- `type` - type of policy.
