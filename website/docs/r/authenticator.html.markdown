---
layout: 'okta'
page_title: 'Okta: okta_authenticator'
sidebar_current: 'docs-okta-resource-okta-authenticator'
description: |-
  Manages Okta Authenticator
---

# okta_authenticator

~> **WARNING:** This feature is only available as a part of the Identity Engine. [Contact support](mailto:dev-inquiries@okta.com) for further information.

This resource allows you to configure different authenticators.

-> **NOTE:** An authenticator can only be deleted if it's not in use by any policy.

## Example Usage

```hcl
resource "okta_authenticator" "test" {
  name     = "Security Question"
  key      = "security_question"
  settings = jsonencode(
  {
    "allowedFor" : "recovery"
  }
  )
}
```

## Argument Reference

The following arguments are supported:

- `key` (Required) A human-readable string that identifies the authenticator. Some authenticators are available by feature flag on the organization. Possible values inclue: `"external_idp"`, `"google_otp"`, `"okta_email"`, `"okta_password"`, `"okta_verify"`, `"onprem_mfa"`, `"phone_number"`, `"rsa_token"`, `"security_question"`, `"webauthn"`, `"duo"`, `"yubikey_token"`.

- `name` - (Required) Name of the authenticator.

- `status` - (Optional) Status of the authenticator. Default is `ACTIVE`.

- `settings` - (Optional) Settings for the authenticator. Settings object contains values based on Authenticator key. It is not used for authenticators with type `"security_key"`.

- `provider_hostname` - (Optional) Server host name or IP address. Default is `"localhost"`. Used only for authenticators with type `"security_key"`.

- `provider_auth_port` - (Optional) The RADIUS server port (for example 1812). This is defined when the On-Prem RADIUS server is configured. Default is `9000`. Used only for authenticators with type `"security_key"`.

- `provider_shared_secret` - (Optional) An authentication key that must be defined when the RADIUS server is configured, and must be the same on both the RADIUS client and server. Used only for authenticators with type `"security_key"`.

- `provider_user_name_template` - (Optional) Username template expected by the provider. Used only for authenticators with type `"security_key"`.

## Attributes Reference

- `id` - ID of the authenticator.

- `type` - Type of the Authenticator.

- `provider_instance_id` - App Instance ID.

- `provider_type` - The type of Authenticator. Values include: `"password"`, `"security_question"`, `"phone"`, `"email"`, `"app"`, `"federated"`, and `"security_key"`.

## Import

Okta authenticator can be imported via the Okta ID.

```
$ terraform import okta_authenticator.example &#60;authenticator_id&#62;
```
