---
page_title: "okta_network_zone Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Network Zone.
---

# okta_network_zone (Data Source)

Use this data source to retrieve an Okta Network Zone.

## Example Usage

```hcl
data "okta_network_zone" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
