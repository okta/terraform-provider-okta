---
page_title: "okta_event_hook Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Event Hook.
---

# okta_event_hook (Data Source)

Use this data source to retrieve an Okta Event Hook.

## Example Usage

```hcl
data "okta_event_hook" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
