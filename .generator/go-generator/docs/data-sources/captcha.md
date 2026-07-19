---
page_title: "okta_captcha Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Captcha.
---

# okta_captcha (Data Source)

Use this data source to retrieve an Okta Captcha.

## Example Usage

```hcl
data "okta_captcha" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
