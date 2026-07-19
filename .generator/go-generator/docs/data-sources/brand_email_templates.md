---
page_title: "okta_brand_email_templates Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Brand Email Templates.
---

# okta_brand_email_templates (Data Source)

Use this data source to retrieve an Okta Brand Email Templates.

## Example Usage

```hcl
data "okta_brand_email_templates" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
