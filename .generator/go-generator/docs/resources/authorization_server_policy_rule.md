---
page_title: "okta_authorization_server_policy_rule Resource - terraform-provider-okta"
description: |-
  Manages an Okta Authorization Server Policy Rule resource.
---

# okta_authorization_server_policy_rule (Resource)

Manages an Okta Authorization Server Policy Rule.

## Example Usage

```hcl
resource "okta_authorization_server_policy_rule" "example" {
  auth_server_id = "<auth-server-id>"
  policy_id = "<policy-id>"

  # Optional fields
  # access_token_lifetime_minutes = 0
  # inline_hook = "<inline_hook>"
  # refresh_token_lifetime_minutes = 0
  # refresh_token_window_minutes = 0
  # include = "<include>"
}
```

## Schema

### Required

- `auth_server_id` (String) ID of the authorization server
- `policy_id` (String) ID of the authorization server policy

### Optional

- `access_token_lifetime_minutes` (Number) Lifetime of the access token in minutes. The minimum is five minutes. The maximum is one day.
- `inline_hook` (String) InlineHook
- `refresh_token_lifetime_minutes` (Number) Lifetime of the refresh token is the minimum access token lifetime.
- `refresh_token_window_minutes` (Number) Timeframe when the refresh token is valid. The minimum is 10 minutes. The maximum is five years (2,628,000 minutes).
- `include` (List) Array of grant types that this condition includes.
- `include` (List) Groups to be included
- `include` (List) Users to be included
- `include` (List) Include
- `name` (String) Name of the rule
- `priority` (Number) Priority of the rule
- `status` (String) Status of the rule
- `system` (Boolean) Set to `true` for system rules. You can't delete system rules.
- `type` (String) Rule type

### Read-Only

- `id` (String) The unique identifier for the resource.
- `created` (String) Timestamp when the rule was created
- `last_updated` (String) Timestamp when the rule was last modified

## Import

Import using `{auth_server_id}/{policy_id}/{id}`:

```shell
terraform import okta_authorization_server_policy_rule.example <auth_server_id> <policy_id> <id>
```
