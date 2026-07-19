---
page_title: "okta_brands Resource - terraform-provider-okta"
description: |-
  Manages an Okta Brands resource.
---

# okta_brands (Resource)

Manages an Okta Brands.

## Example Usage

```hcl
resource "okta_brands" "example" {

  # Optional fields
  # agree_to_custom_privacy_policy = true
  # custom_privacy_policy_url = "https://example.com"
  # app_instance_id = "<app-instance-id>"
  # app_link_name = "Example App Link Name"
  # classic_application_uri = "https://example.com"
}
```

## Schema

### Optional

- `agree_to_custom_privacy_policy` (Boolean) Consent for updating the custom privacy URL. Not required when resetting the URL.
- `custom_privacy_policy_url` (String) Custom privacy policy URL
- `app_instance_id` (String) ID for the App instance
- `app_link_name` (String) Name for the app instance
- `classic_application_uri` (String) Application URI for classic Orgs
- `email_domain_id` (String) The ID of the email domain
- `locale` (String) The language specified as an [IETF BCP 47 language tag](https://datatracker.ietf.org/doc/html/rfc5646)
- `name` (String) The name of the Brand
- `remove_powered_by_okta` (Boolean) Removes 'Powered by Okta' from the sign-in page in redirect authentication deployments, and '© [current year] Okta, Inc.' from the Okta End-User Dashboard

### Read-Only

- `id` (String) The unique identifier for the resource.
- `embedded` (String) Embedded
- `is_default` (Boolean) If `true`, the Brand is used for the Okta subdomain
