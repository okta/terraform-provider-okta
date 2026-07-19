---
page_title: "okta_email_server Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Email Server.
---

# okta_email_server (Data Source)

Use this data source to retrieve an Okta Email Server.

## Example Usage

```hcl
data "okta_email_server" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
