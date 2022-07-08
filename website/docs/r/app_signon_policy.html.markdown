---
layout: "okta"
page_title: "Okta: okta_app_signon_policy"
sidebar_current: "docs-okta-resource-okta-app-signon-policy"
description: |-
  Manages a sign-on policy.
---

# okta_app_signon_policy

~> **WARNING:** This feature is only available as a part of the Identity Engine. [Contact support](mailto:dev-inquiries@okta.com) for further information.

This resource allows you to create and configure a sign-on policy for the application. (Inside the product this is referenced as an _authentication policy_)

A newly create app sign-on policy will always contain a default `Catch-all Rule`.

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

\_The same mechanism is in place for `okta_app_oauth` and `okta_app_saml`.

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
