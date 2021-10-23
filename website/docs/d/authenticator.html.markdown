---
layout: 'okta'
page_title: 'Okta: okta_authenticator'
sidebar_current: 'docs-okta-datasource-okta-authenticator'
description: |-
  Manages Okta Authenticator
---

# okta_authenticator

~> **WARNING:** This feature is only available as a part of the Identity Engine. [Contact support](mailto:dev-inquiries@okta.com) for further information.

Use this data source to retrieve an authenticator.

## Example Usage

```hcl
data "okta_authenticator" "test" {
  name = "Security Question"
}
```

```hcl
data "okta_authenticator" "test" {
  key = "okta_email"
}
```

## Argument Reference

The following arguments are supported:

- `key` (Optional) A human-readable string that identifies the authenticator.

- `name` - (Optional) Name of the authenticator.

- `id` - (Optional) ID of the authenticator.

## Attributes Reference

- `id` - ID of the authenticator.

- `type` - Type of the Authenticator.

- `name` - Name of the authenticator.

- `type` - Type of the Authenticator.

- `settings` - Settings for the authenticator.


