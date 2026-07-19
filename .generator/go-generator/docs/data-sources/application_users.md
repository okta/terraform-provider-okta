---
page_title: "okta_application_users Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Application Users.
---

# okta_application_users (Data Source)

Use this data source to retrieve an Okta Application Users.

## Example Usage

```hcl
data "okta_application_users" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
