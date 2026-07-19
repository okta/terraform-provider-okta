---
page_title: "okta_email_domain Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Email Domain.
---

# okta_email_domain (Data Source)

Use this data source to retrieve an Okta Email Domain.

## Example Usage

```hcl
data "okta_email_domain" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
