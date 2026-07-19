---
page_title: "okta_application_group_assignments Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Application Group Assignments.
---

# okta_application_group_assignments (Data Source)

Use this data source to retrieve an Okta Application Group Assignments.

## Example Usage

```hcl
data "okta_application_group_assignments" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
