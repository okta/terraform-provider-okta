---
page_title: "okta_ssf_receiver Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Ssf Receiver.
---

# okta_ssf_receiver (Data Source)

Use this data source to retrieve an Okta Ssf Receiver.

## Example Usage

```hcl
data "okta_ssf_receiver" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
