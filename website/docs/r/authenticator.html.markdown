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

-> **Create:** The Okta API has an odd notion of create for authenticators. If
the authenticator doesn't exist then a one time `POST /api/v1/authenticators` to
create the authenticator (hard create) will be performed. Thereafter, that
authenticator is never deleted, it is only deactivated (soft delete). Therefore,
if the authenticator already exists create is just a soft import of an existing
authenticator.

-> **Delete:** Authenticators can not be truly deleted therefore delete is soft.
Delete will attempt to deativate the authenticator. An authenticator can only be
deactivated if it's not in use by any other policy.

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

- `key` (Required) A human-readable string that identifies the authenticator. Some authenticators are available by feature flag on the organization. Possible values inclue: `duo`, `external_idp`, `google_otp`, `okta_email`, `okta_password`, `okta_verify`, `onprem_mfa`, `phone_number`, `rsa_token`, `security_question`, `webauthn`

- `name` - (Required) Name of the authenticator.

- `status` - (Optional) Status of the authenticator. Default is `ACTIVE`.

- `settings` - (Optional) Settings for the authenticator. The settings JSON contains values based on Authenticator key. It is not used for authenticators with type `"security_key"`.

- `provider_json` - (Optional) Provider JSON allows for expressive provider
values. This argument conflicts with the other `provider_xxx` arguments.  The
[Create
Provider](https://developer.okta.com/docs/reference/api/authenticators-admin/#request)
illustrates detailed provider values for a Duo authenticator.  [Provider
values](https://developer.okta.com/docs/reference/api/authenticators-admin/#authenticators-administration-api-object)
are listed in Okta API.
- `provider_auth_port` - (Optional) The RADIUS server port (for example 1812). This is defined when the On-Prem RADIUS server is configured. Used only for authenticators with type `"security_key"`.  Conflicts with `provider_json` argument.

- `provider_hostname` - (Optional) Server host name or IP address. Default is `"localhost"`. Used only for authenticators with type `"security_key"`.  Conflicts with `provider_json` argument.


- `provider_shared_secret` - (Optional) An authentication key that must be defined when the RADIUS server is configured, and must be the same on both the RADIUS client and server. Used only for authenticators with type `"security_key"`.  Conflicts with `provider_json` argument.

- `provider_user_name_template` - (Optional) Username template expected by the provider. Used only for authenticators with type `"security_key"`.  Conflicts with `provider_json` argument.

- `provider_host` - (Optional) (DUO specific) - The Duo Security API hostname". Conflicts with `provider_json` argument.

- `provider_integration_key` - (Optional) (DUO specific) - The Duo Security integration key.  Conflicts with `provider_json` argument.

- `provider_secret_key` - (Optional) (DUO specific) - The Duo Security secret key.  Conflicts with `provider_json` argument.

## Attributes Reference

- `id` - ID of the authenticator.

- `type` - The type of Authenticator. Values include: `"password"`, `"security_question"`, `"phone"`, `"email"`, `"app"`, `"federated"`, and `"security_key"`.

- `provider_instance_id` - App Instance ID.

- `provider_type` - Provider type. Supported value for Duo: `DUO`. Supported value for Custom App: `PUSH`

## Import

Okta authenticator can be imported via the Okta ID.

```
$ terraform import okta_authenticator.example &#60;authenticator_id&#62;
```
