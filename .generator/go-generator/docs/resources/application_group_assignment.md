---
page_title: "okta_application_group_assignment Resource - terraform-provider-okta"
description: |-
  Manages an Okta Application Group Assignment resource.
---

# okta_application_group_assignment (Resource)

Manages an Okta Application Group Assignment.

## Example Usage

```hcl
resource "okta_application_group_assignment" "example" {
  app_id = "<app-id>"

  # Optional fields
  # last_updated = "<last_updated>"
  # priority = 0
  # profile = "<profile>"
}
```

## Schema

### Required

- `app_id` (String) ID of the parent application

### Optional

- `last_updated` (String) LastUpdated
- `priority` (Number) Priority assigned to the group. If an app has more than one group assigned to the same user, then the group with the higher priority has its profile applied to the [application user](https://develo...
- `profile` (String) Specifies the profile properties applied to [application users](https://developer.okta.com/docs/api/openapi/okta-management/management/tags/applicationusers) that are assigned to the app through gr...

### Read-Only

- `id` (String) The unique identifier for the resource.
- `embedded` (String) Embedded resource related to the Application Group using the [JSON Hypertext Application Language](https://datatracker.ietf.org/doc/html/draft-kelly-json-hal-06) specification. If the `expand=group...

## Import

Import using `{app_id}/{id}`:

```shell
terraform import okta_application_group_assignment.example <app_id> <id>
```
