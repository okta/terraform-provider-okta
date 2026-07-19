---
page_title: "okta_identity_provider_signing_keys Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Identity Provider Signing Keys.
---

# okta_identity_provider_signing_keys (Data Source)

Use this data source to retrieve an Okta Identity Provider Signing Keys.

## Example Usage

```hcl
data "okta_identity_provider_signing_keys" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
