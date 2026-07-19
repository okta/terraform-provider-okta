---
page_title: "okta_linked_object Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Linked Object.
---

# okta_linked_object (Data Source)

Use this data source to retrieve an Okta Linked Object.

## Example Usage

```hcl
data "okta_linked_object" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
