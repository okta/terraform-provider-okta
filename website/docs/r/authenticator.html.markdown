---
layout: 'okta'
page_title: 'Okta: okta_authenticator'
sidebar_current: 'docs-okta-resource-authenticator'
description: |-
  Updates an Okta Authenticator.
---

# okta_authenticator data

Data is supported actinging a find for the corresponding authenticator.

```hcl
data "okta_authenticator" "test" {
  type = "security_question"
}
```

# okta_authenticator resource

Updates an Okta Authenticator.

This resource allows you to update the status on an Okta Authenticator

## Example Usage

```hcl
resource "okta_authenticator" "example" {
  type = "security_question"
  status = "ACTIVE"
}
```

## Argument Reference

The following arguments are supported:

- `status` - The authenticator's status `ACTIVE` | `INACTIVE`

## Attributes Reference

- `id` - The ID of the Okta Authenticator.
- `key` - The key value for the Okta Authenticator.
- `key` - The name value for the Okta Authenticator.
- `settings` - The settings of the Okta Authenticator in JSON format.
- `status` - The status of the Okta Authenticator `ACTIVE` | `INACTIVE`.
- `type` - The type value for the Okta Authenticator.

## Import

Import is not allowed. Existence of authenticator resources is immutable therefore deletion and creation is not allowed.

