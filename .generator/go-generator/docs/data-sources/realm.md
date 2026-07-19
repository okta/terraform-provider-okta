---
page_title: "okta_realm Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Realm.
---

# okta_realm (Data Source)

Use this data source to retrieve an Okta Realm.

## Example Usage

```hcl
data "okta_realm" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
