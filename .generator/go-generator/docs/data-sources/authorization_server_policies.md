---
page_title: "okta_authorization_server_policies Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Authorization Server Policies.
---

# okta_authorization_server_policies (Data Source)

Use this data source to retrieve an Okta Authorization Server Policies.

## Example Usage

```hcl
data "okta_authorization_server_policies" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
