---
layout: 'okta'
page_title: 'Okta: okta_policy_rule_idp_discovery'
sidebar_current: 'docs-okta-resource-policy-rule-idp-discovery'
description: |-
  Creates an IdP Discovery Policy Rule.
---

# okta_policy_rule_idp_discovery

Creates an IdP Discovery Policy Rule.

This resource allows you to create and configure an IdP Discovery Policy Rule.

## Example Usage

```hcl
resource "okta_policy_rule_idp_discovery" "example" {
  policyid                  = "<policy id>"
  name                      = "example"
  idp_id                    = "<idp id>"
  idp_type                  = "OIDC"
  network_connection        = "ANYWHERE"
  priority                  = 1
  status                    = "ACTIVE"
  user_identifier_type      = "ATTRIBUTE"
  user_identifier_attribute = "company"

  app_exclude {
    id   = "<app id>"
    type = "APP"
  }

  app_exclude {
    name = "yahoo_mail"
    type = "APP_TYPE"
  }

  app_include {
    id   = "<app id>"
    type = "APP"
  }

  app_include {
    name = "<app type name>"
    type = "APP_TYPE"
  }

  platform_include {
    type    = "MOBILE"
    os_type = "OSX"
  }

  user_identifier_patterns {
    match_type = "EQUALS"
    value      = "Articulate"
  }
}
```

## Argument Reference

The following arguments are supported:

- `policyid` - (Required) Policy ID.

- `name` - (Required) Policy rule name.

- `idp_id` - (Optional) The identifier for the Idp the rule should route to if all conditions are met.

- `idp_type` - (Optional) Type of Idp. One of: `"SAML2"`, `"IWA"`, `"AgentlessDSSO"`, `"X509"`, `"FACEBOOK"`, `"GOOGLE"`, `"LINKEDIN"`, `"MICROSOFT"`, `"OIDC"`

- `network_connection` - (Optional) The network selection mode. One of `"ANYWEHRE"` or `"ZONE"`.

- `network_includes` - Required if `network_connection` = `"ZONE"`. Indicates the network zones to include.

- `network_excludes` - Required if `network_connection` = `"ZONE"`. Indicates the network zones to exclude.

- `priority` - (Optional) Idp rule priority. This attribute can be set to a valid priority. To avoid an endless diff situation an error is thrown if an invalid property is provided. The Okta API defaults to the last/lowest if not provided.

- `status` - (Optional) Idp rule status: `"ACTIVE"` or `"INACTIVE"`. By default it is `"ACTIVE"`.

- `user_identifier_type` - (Optional) One of: `"IDENTIFIER"`, `"ATTRIBUTE"`

- `user_identifier_attribute` - (Optional) Profile attribute matching can only have a single value that describes the type indicated in `user_identifier_type`. This is the attribute or identifier that the `user_identifier_patterns` are checked against.

- `app_include` - (Optional) Applications to include in discovery rule.

  - `id` - (Optional) Use if `type` is `"APP"` to indicate the application Id to include.

  - `name` - (Optional) Use if the `type` is `"APP_TYPE"` to indicate the type of application(s) to include in instances where an entire group (i.e. `yahoo_mail`) of applications should be included.

  - `type` - (Optional) One of: `"APP"`, `"APP_TYPE"`

```hcl
app_include {
  id = string
  type = string
  name = string
}
```

- `app_exclude` - (Optional) Applications to exclude in discovery. See `app_include` for details.

```hcl
app_exclude {
  id = string
  type = string
  name = string
}
```

- `platform_include` - (Optional)

  - `type` - (Optional) One of: `"ANY"`, `"MOBILE"`, `"DESKTOP"`

  - `os_expression` - (Optional) Only available when using `os_type = "OTHER"`

  - `os_type` - (Optional) One of: `"ANY"`, `"IOS"`, `"WINDOWS"`, `"ANDROID"`, `"OTHER"`, `"OSX"`

```hcl
app_exclude {
  type = string
  os_expression = string
  os_type = string
}
```

- `user_identifier_patterns` - (Optional) Specifies a User Identifier pattern condition to match against. If `match_type` of `"EXPRESSION"` is used, only a *single* element can be set. Otherwise multiple elements of matching patterns may be provided.

  - `match_type` - (Optional) The kind of pattern. For regex, use `"EXPRESSION"`. For simple string matches, use one of the following: `"SUFFIX"`, `"EQUALS"`, `"STARTS_WITH"`, `"CONTAINS"`

  - `value` - (Optional) The regex or simple match string to match against.

```hcl
user_identifier_patterns {
  match_type = string
  value = string
}
```

## Attributes Reference

- `id` - ID of the Rule.

- `policyid` - Policy ID.

## Import

A Policy Rule can be imported via the Policy and Rule ID.

```
$ terraform import okta_policy_rule_idp_discovery.example <policy id>/<rule id>
```
