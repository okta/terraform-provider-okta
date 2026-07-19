---
page_title: "okta_device_assurance Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Device Assurance.
---

# okta_device_assurance (Data Source)

Use this data source to retrieve an Okta Device Assurance.

## Example Usage

```hcl
data "okta_device_assurance" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
