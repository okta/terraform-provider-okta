---
layout: 'okta'
page_title: 'Okta: okta_policy_rule_mfa'
sidebar_current: 'docs-okta-resource-policy-rule-mfa'
description: |-
  Creates an MFA Policy Rule.
---

# okta_policy_rule_mfa

This resource allows you to create and configure an MFA Policy Rule.

## Example Usage

```hcl
data "okta_default_policy" "example" {
  type = "MFA_ENROLL"
}

resource "okta_policy_rule_mfa" "example" {
  policy_id = data.okta_default_policy.example.id
  name      = "My Rule"
  status    = "ACTIVE"
  enroll    = "LOGIN"
  app_include {
    id   = okta_app_oauth.example.id
    type = "APP"
  }
  app_include {
    type = "APP_TYPE"
    name = "yahoo_mail"
  }
}

resource "okta_app_oauth" "example" {
  label          = "My App"
  type           = "web"
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://localhost:8000"]
  response_types = ["code"]
  skip_groups    = true
}
```

## Argument Reference

The following arguments are supported:

- `policyid` - (Deprecated) Policy ID.
  
- `policy_id` - (Required) Policy ID.

- `name` - (Required) Policy Rule Name.

- `priority` - (Optional) Policy Rule Priority, this attribute can be set to a valid priority. To avoid endless diff situation we error if an invalid priority is provided. API defaults it to the last (lowest) if not there.

- `status` - (Optional) Policy Rule Status: `"ACTIVE"` or `"INACTIVE"`.

- `enroll` - (Optional) When a user should be prompted for MFA. It can be `"CHALLENGE"`, `"LOGIN"`, or `"NEVER"`.

- `network_connection` - (Optional) Network selection mode: `"ANYWHERE"`, `"ZONE"`, `"ON_NETWORK"`, or `"OFF_NETWORK"`.

- `network_includes` - (Optional) The network zones to include. Conflicts with `network_excludes`.

- `network_excludes` - (Optional) The network zones to exclude. Conflicts with `network_includes`.

- `app_include` - (Optional) Applications to include in discovery rule. **IMPORTANT**: this field is only available in Classic Organizations.
  - `id` - (Optional) Use if `type` is `"APP"` to indicate the application id to include.
  - `name` - (Optional) Use if the `type` is `"APP_TYPE"` to indicate the type of application(s) to include in instances where an entire group (i.e. `yahoo_mail`) of applications should be included.
  - `type` - (Required) One of: `"APP"`, `"APP_TYPE"`

## Attributes Reference

- `id` - ID of the Rule.

- `policyid` - (Deprecated) Policy ID.
  
- `policy_id` - Policy ID.

## Import

A Policy Rule can be imported via the Policy and Rule ID.

```
$ terraform import okta_policy_rule_mfa.example <policy id>/<rule id>
```
