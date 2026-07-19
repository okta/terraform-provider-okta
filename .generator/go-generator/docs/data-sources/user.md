---
page_title: "okta_user Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta User.
---

# okta_user (Data Source)

Use this data source to retrieve an Okta User.

## Example Usage

```hcl
data "okta_user" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
