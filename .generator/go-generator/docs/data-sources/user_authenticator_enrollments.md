---
page_title: "okta_user_authenticator_enrollments Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta User Authenticator Enrollments.
---

# okta_user_authenticator_enrollments (Data Source)

Use this data source to retrieve an Okta User Authenticator Enrollments.

## Example Usage

```hcl
data "okta_user_authenticator_enrollments" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
