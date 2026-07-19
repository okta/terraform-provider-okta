---
page_title: "okta_application_tokens Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Application Tokens.
---

# okta_application_tokens (Data Source)

Use this data source to retrieve an Okta Application Tokens.

## Example Usage

```hcl
data "okta_application_tokens" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
