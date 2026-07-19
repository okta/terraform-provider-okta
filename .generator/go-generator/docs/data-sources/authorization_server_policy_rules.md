---
page_title: "okta_authorization_server_policy_rules Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Authorization Server Policy Rules.
---

# okta_authorization_server_policy_rules (Data Source)

Use this data source to retrieve an Okta Authorization Server Policy Rules.

## Example Usage

```hcl
data "okta_authorization_server_policy_rules" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
