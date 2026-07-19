---
page_title: "okta_captcha Resource - terraform-provider-okta"
description: |-
  Manages an Okta Captcha resource.
---

# okta_captcha (Resource)

Manages an Okta Captcha.

## Example Usage

```hcl
resource "okta_captcha" "example" {

  # Optional fields
  # name = "Example Name"
  # secret_key = "<secret_key>"
  # site_key = "<site_key>"
  # type = "<type>"
}
```

## Schema

### Optional

- `name` (String) The name of the CAPTCHA instance
- `secret_key` (String) The secret key issued from the CAPTCHA provider to perform server-side validation for a CAPTCHA token
- `site_key` (String) The site key issued from the CAPTCHA provider to render a CAPTCHA on a page
- `type` (String) The type of CAPTCHA provider

### Read-Only

- `id` (String) The unique identifier for the resource.
