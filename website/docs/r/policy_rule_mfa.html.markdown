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

Unchecked `Okta` and checked `Applications` (with `Any application that supports MFA enrollment` option) checkboxes in the `User is accessing` section corresponds to the following config:

```hcl
data "okta_default_policy" "example" {
  type = "MFA_ENROLL"
}

resource "okta_policy_rule_mfa" "example" {
  name      = "Some policy rule"
  policy_id = data.okta_default_policy.example.id

  app_exclude {
    name = "okta"
    type = "APP_TYPE"
  }
}
```

Unchecked `Okta` and checked `Applications` (with `Specific applications` option) checkboxes in the `User is accessing` section corresponds to the following config:

```hcl
data "okta_default_policy" "example" {
  type = "MFA_ENROLL"
}

resource "okta_policy_rule_mfa" "example" {
  name      = "Some policy rule"
  policy_id = data.okta_default_policy.example.id

  app_exclude {
    name = "okta"
    type = "APP_TYPE"
  }

  app_include {
    id   = "some_app_id"
    type = "APP"
  }
}
```

Checked `Okta` and unchecked `Applications` checkboxes in the `User is accessing` section corresponds to the following config:

```hcl
data "okta_default_policy" "example" {
  type = "MFA_ENROLL"
}

resource "okta_policy_rule_mfa" "example" {
  name      = "Some policy rule"
  policy_id = data.okta_default_policy.example.id

  app_include {
    name = "okta"
    type = "APP_TYPE"
  }
}
```

Checked `Okta` and checked `Applications` (with `Any application that supports MFA enrollment` option) checkboxes in the `User is accessing` section corresponds to the following config:

```hcl
data "okta_default_policy" "example" {
  type = "MFA_ENROLL"
}

resource "okta_policy_rule_mfa" "example" {
  name      = "Some policy rule"
  policy_id = data.okta_default_policy.example.id
}
```

Checked `Okta` and checked `Applications` (with `Specific applications` option) checkboxes in the `User is accessing` section corresponds to the following config:

```hcl
data "okta_default_policy" "example" {
  type = "MFA_ENROLL"
}

resource "okta_policy_rule_mfa" "example" {
  name      = "Some policy rule"
  policy_id = data.okta_default_policy.example.id

  app_include {
    name = "okta"
    type = "APP_TYPE"
  }

  app_include {
    id   = "some_app_id"
    type = "APP"
  }
}
```

## Argument Reference

The following arguments are supported:
  
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

- `app_exlcude` - (Optional) Applications to exclude from the discovery rule. **IMPORTANT**: this field is only available in Classic Organizations.
  - `id` - (Optional) Use if `type` is `"APP"` to indicate the application id to include.
  - `name` - (Optional) Use if the `type` is `"APP_TYPE"` to indicate the type of application(s) to include in instances where an entire group (i.e. `yahoo_mail`) of applications should be included.
  - `type` - (Required) One of: `"APP"`, `"APP_TYPE"`

## Attributes Reference

- `id` - ID of the Rule.
  
- `policy_id` - Policy ID.

## Import

A Policy Rule can be imported via the Policy and Rule ID.

```
$ terraform import okta_policy_rule_mfa.example &#60;policy id&#62;/&#60;rule id&#62;
```
