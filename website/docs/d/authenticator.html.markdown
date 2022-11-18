---
layout: 'okta'
page_title: 'Okta: okta_authenticator'
sidebar_current: 'docs-okta-datasource-okta-authenticator'
description: |-
  Get an authenticator by key, name of ID.
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

- `id` - (Optional) ID of the authenticator.

- `key` (Optional) A human-readable string that identifies the authenticator.

- `name` - (Optional) Name of the authenticator.

## Attributes Reference

- `id` - ID of the authenticator.

- `name` - Name of the authenticator.

- `provider_auth_port` - (Specific to `security_key`) The provider server port (for example 1812).

- `provider_hostname` - (Specific to `security_key`) Server host name or IP address.

- `provider_instance_id` - (Specific to `security_key`) App Instance ID.

- `provider_type` - Provider type.

- `provider_user_name_template` - Username template expected by the provider.

- `provider_host` - (Specific to `DUO` provider) The Duo Security API hostname.

- `provider_integration_key` - (Specific to `DUO` provider) The Duo Security integration key.

- `provider_secret_key` - (Specific to `DUO` provider) The Duo Security secret key.

- `settings` - Settings for the authenticator (expressed in JSON).

- `status` - Status of the Authenticator.

- `type` - The type of Authenticator.
