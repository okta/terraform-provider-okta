---
layout: "okta"
page_title: "Okta: okta_factor"
sidebar_current: "docs-okta-resource-factor"
description: |-
  Allows you to manage the activation of Okta MFA methods.
---

# okta_factor

Allows you to manage the activation of Okta MFA methods.

This resource allows you to manage Okta MFA methods.

## Example Usage

```hcl
resource "okta_factor" "example" {
  provider_id = "google_otp"
}
```

## Argument Reference

The following arguments are supported:

* `provider_id` - (Required) The MFA provider name.

* `active` - (Optional) Whether or not to activate the provider, by default it is set to `true`.

## Attributes Reference

* `provider_id` - MFA provider name.
