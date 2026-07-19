---
page_title: "okta_authorization_server Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Authorization Server.
---

# okta_authorization_server (Data Source)

Use this data source to retrieve an Okta Authorization Server.

## Example Usage

```hcl
data "okta_authorization_server" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
