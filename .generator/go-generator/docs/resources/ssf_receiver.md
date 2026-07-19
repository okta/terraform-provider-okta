---
page_title: "okta_ssf_receiver Resource - terraform-provider-okta"
description: |-
  Manages an Okta Ssf Receiver resource.
---

# okta_ssf_receiver (Resource)

Manages an Okta Ssf Receiver.

## Example Usage

```hcl
resource "okta_ssf_receiver" "example" {

  # Optional fields
  # name = "Example Name"
  # issuer = "<issuer>"
  # jwks_url = "https://example.com"
  # well_known_url = "https://example.com"
  # type = "<type>"
}
```

## Schema

### Optional

- `name` (String) The name of the security events provider instance
- `issuer` (String) Issuer URL
- `jwks_url` (String) The public URL where the JWKS public key is uploaded
- `well_known_url` (String) The well-known URL of the security events provider (the SSF transmitter)
- `type` (String) The app type of the security events provider

### Read-Only

- `id` (String) The unique identifier for the resource.
- `status` (String) Indicates whether the security events provider is active or not
