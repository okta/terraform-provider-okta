---
page_title: "okta_identity_source_sessions Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Identity Source Sessions.
---

# okta_identity_source_sessions (Data Source)

Use this data source to retrieve an Okta Identity Source Sessions.

## Example Usage

```hcl
data "okta_identity_source_sessions" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
