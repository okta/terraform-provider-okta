---
page_title: "okta_authorization_server_claims Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Authorization Server Claims.
---

# okta_authorization_server_claims (Data Source)

Use this data source to retrieve an Okta Authorization Server Claims.

## Example Usage

```hcl
data "okta_authorization_server_claims" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
