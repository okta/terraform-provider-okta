---
page_title: "okta_authenticator Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Authenticator.
---

# okta_authenticator (Data Source)

Use this data source to retrieve an Okta Authenticator.

## Example Usage

```hcl
data "okta_authenticator" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
