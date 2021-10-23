---
layout: 'okta' 
page_title: 'Okta: okta_app_signon_policy' 
sidebar_current: 'docs-okta-datasource-app-signon-policy'
description: |- 
    Get a sign-on policy for the application.
---

# okta_app_signon_policy

~> **WARNING:** This feature is only available as a part of the Identity Engine. [Contact support](mailto:dev-inquiries@okta.com) for further information.

Use this data source to retrieve a sign-on policy for the application.

## Example Usage

```hcl
data "okta_app_signon_policy" "example" {
  app_id = "app_id"
}
```

## Arguments Reference

- `app_id` - (Required) The application ID.

## Attributes Reference

- `id` - Sign-on policy ID.
