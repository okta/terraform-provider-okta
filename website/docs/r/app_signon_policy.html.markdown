---
layout: "okta"
page_title: "Okta: okta_app_signon_policy"
sidebar_current: "docs-okta-resource-okta-app-signon-policy"
description: |-
  Manages a sign-on policy.
---

# okta_app_signon_policy

~> **WARNING:** This feature is only available as a part of the Okta Identity
Engine (OIE) and ***is not*** compatible with Classic orgs. Authentication
policies for applications in a Classic org can only be modified in the Admin UI,
there isn't a public API for this. Therefore the Okta Terraform Provider does
not support this resource for Classic orgs. [Contact
support](mailto:dev-inquiries@okta.com) for further information.

This resource allows you to create and configure a sign-on policy for the
application. Inside the product a sign-on policy is referenced as an
_authentication policy_, in the public API the policy is of type
[`ACCESS_POLICY`](https://developer.okta.com/docs/reference/api/policy/#policy-object).

A newly created app's sign-on policy will always contain the default
authentication policy unless one is assigned via `authentication_policy` in the
app resource. At the API level the default policy has system property value of
true.

~> **WARNING:** When this policy is destroyed any other applications that
associate the policy as their authentication policy will be reassigned to the
default/system access policy.

## Example Usage

```hcl
resource "okta_app_oauth" "my_app" {
  label                     = "My App"
  type                      = "web"
  grant_types               = ["authorization_code"]
  redirect_uris             = ["http://localhost:3000"]
  post_logout_redirect_uris = ["http://localhost:3000"]
  response_types            = ["code"]
  // this is needed to associate the application with the policy
  authentication_policy     = okta_app_signon_policy.my_app_policy.id
}

resource "okta_app_signon_policy" "my_app_policy" {
  name        = "My App Sign-On Policy"
  description = "Authentication Policy to be used on my app."
}
```

The created policy can be extended using `app_signon_policy_rules`.

```hcl
resource "okta_app_signon_policy" "my_app_policy" {
  name        = "My App Sign-On Policy"
  description = "Authentication Policy to be used on my app."
}

resource "okta_app_signon_policy_rule" "some_rule" {
  policy_id                   = resource.okta_app_signon_policy.my_app_policy.id
  name                        = "Some Rule"
  factor_mode                 = "1FA"
  re_authentication_frequency = "PT43800H"
  constraints = [
    jsonencode({
      "knowledge" : {
        "types" : ["password"]
      }
    })
  ]
}
```

## Argument Reference

- `name` - (Required) Name of the policy.
- `description` - (Required) Description of the policy.

## Attributes Reference

- `id` - ID of the sign-on policy.
