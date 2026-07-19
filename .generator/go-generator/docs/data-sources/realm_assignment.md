---
page_title: "okta_realm_assignment Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Realm Assignment.
---

# okta_realm_assignment (Data Source)

Use this data source to retrieve an Okta Realm Assignment.

## Example Usage

```hcl
data "okta_realm_assignment" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
