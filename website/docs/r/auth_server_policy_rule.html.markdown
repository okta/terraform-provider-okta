---
layout: 'okta'
page_title: 'Okta: okta_auth_server_policy_rule'
sidebar_current: 'docs-okta-resource-auth-server-policy-rule'
description: |-
  Creates an Authorization Server Policy Rule.
---

# okta_auth_server_policy_rule

Creates an Authorization Server Policy Rule.

This resource allows you to create and configure an Authorization Server Policy Rule.

-> This resource is concurrency safe. However, when creating/updating/deleting
multiple rules belonging to a policy, the Terraform meta argument
[`depends_on`](https://www.terraform.io/language/meta-arguments/depends_on)
should be added to each rule chaining them all in sequence. Base the sequence on
the `priority` property in ascending value.

## Example Usage

```hcl
resource "okta_auth_server_policy_rule" "example" {
  auth_server_id       = "<auth server id>"
  policy_id            = "<auth server policy id>"
  status               = "ACTIVE"
  name                 = "example"
  priority             = 1
  group_whitelist      = ["<group ids>"]
  grant_type_whitelist = ["implicit"]
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) Auth Server Policy Rule name.

- `auth_server_id` - (Required) Auth Server ID.

- `policy_id` - (Required) Auth Server Policy ID.

- `status` - (Optional) The status of the Auth Server Policy Rule.

- `priority` - (Required) Priority of the auth server policy rule.

- `user_whitelist` - (Optional) Specifies a set of Users to be included.

- `user_blacklist` - (Optional) Specifies a set of Users to be excluded.

- `group_whitelist` - (Optional) Specifies a set of Groups whose Users are to be included. Can be set to Group ID or to the following: `"EVERYONE"`.

- `group_blacklist` - (Optional) Specifies a set of Groups whose Users are to be excluded.

- `grant_type_whitelist` - (Required) Accepted grant type values, `"authorization_code"`, `"implicit"`, `"password"`, `"client_credentials"`, 
  `"urn:ietf:params:oauth:grant-type:saml2-bearer"` (*Early Access Property*), `"urn:ietf:params:oauth:grant-type:token-exchange"` (*Early Access Property*),
  `"urn:ietf:params:oauth:grant-type:device_code"` (*Early Access Property*), `"interaction_code"` (*OIE only*). For `"implicit"` value either `user_whitelist` or `group_whitelist` should be set.

- `scope_whitelist` - (Required) Scopes allowed for this policy rule. They can be whitelisted by name or all can be whitelisted with `"*"`.

- `access_token_lifetime_minutes` - (Optional) Lifetime of access token. Can be set to a value between 5 and 1440 minutes. Default is `60`.

- `refresh_token_lifetime_minutes` - (Optional) Lifetime of refresh token.

- `refresh_token_window_minutes` - (Optional) Window in which a refresh token can be used. It can be a value between 5 and 2628000 (5 years) minutes. Default is `10080` (7 days).
  `"refresh_token_window_minutes"` must be between `"access_token_lifetime_minutes"` and `"refresh_token_lifetime_minutes"`.

- `inline_hook_id` - (Optional) The ID of the inline token to trigger.

## Attributes Reference

- `id` - (Required) The ID of the Auth Server Policy Rule.

- `policy_id` - (Required) The ID of the Auth Server Policy.

- `auth_server_id` - (Required) The ID of the Auth Server.

- `type` - The type of the Auth Server Policy Rule.

## Import

Authorization Server Policy Rule can be imported via the Auth Server ID, Policy ID, and Policy Rule ID.

```
$ terraform import okta_auth_server_policy_rule.example &#60;auth server id&#62;/&#60;policy id&#62;/&#60;policy rule id&#62;
```
