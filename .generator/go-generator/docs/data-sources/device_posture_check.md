---
page_title: "okta_device_posture_check Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Device Posture Check.
---

# okta_device_posture_check (Data Source)

Use this data source to retrieve an Okta Device Posture Check.

## Example Usage

```hcl
data "okta_device_posture_check" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
