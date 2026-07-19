---
page_title: "okta_trusted_origin Resource - terraform-provider-okta"
description: |-
  Manages an Okta Trusted Origin resource.
---

# okta_trusted_origin (Resource)

Manages an Okta Trusted Origin.

## Example Usage

```hcl
resource "okta_trusted_origin" "example" {

  # Optional fields
  # created_by = "<created_by>"
  # last_updated_by = "<last_updated_by>"
  # name = "Example Name"
  # origin = "<origin>"
  # allowed_okta_apps = "<allowed_okta_apps>"
}
```

## Schema

### Optional

- `created_by` (String) The ID of the user who created the trusted origin
- `last_updated_by` (String) The ID of the user who last updated the trusted origin
- `name` (String) Unique name for the trusted origin
- `origin` (String) Unique origin URL for the trusted origin. The supported schemes for this attribute are HTTP, HTTPS, FTP, Ionic 2, and Capacitor.
- `allowed_okta_apps` (List) The allowed Okta apps for the trusted origin scope
- `type` (String) The scope type. Supported values: When you use `IFRAME_EMBED` as the scope type, leave the `allowedOktaApps` property empty to allow iFrame embedding of only Okta sign-in pages. Include `OKTA_ENDUS...
- `status` (String) Status

### Read-Only

- `id` (String) The unique identifier for the resource.
- `created` (String) Timestamp when the trusted origin was created
- `last_updated` (String) Timestamp when the trusted origin was last updated
