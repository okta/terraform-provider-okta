---
page_title: "okta_application_user Resource - terraform-provider-okta"
description: |-
  Manages an Okta Application User resource.
---

# okta_application_user (Resource)

Manages an Okta Application User.

## Example Usage

```hcl
resource "okta_application_user" "example" {
  app_id = "<app-id>"

  # Optional fields
  # created = "<created>"
  # value = "<value>"
  # user_name = "Example User Name"
  # last_updated = "<last_updated>"
  # profile = "<profile>"
}
```

## Schema

### Required

- `app_id` (String) ID of the parent application

### Optional

- `created` (String) Created
- `value` (String) Password value
- `user_name` (String) The user's username in the app  > **Note:** The [userNameTemplate](/openapi/okta-management/management/tags/application/other/createapplication#application/createapplication/t=response&c=200&path=&...
- `last_updated` (String) LastUpdated
- `profile` (String) Specifies the default and custom profile properties for a user. Properties that are visible in the Admin Console for an app assignment can also be assigned through the API. Some properties are refe...
- `scope` (String) Indicates if the assignment is direct (`USER`) or by group membership (`GROUP`). If not specified, Okta tries to determine the scope based on the assignment type.

### Read-Only

- `id` (String) The unique identifier for the resource.
- `embedded` (String) Embedded resources related to the application user using the [JSON Hypertext Application Language](https://datatracker.ietf.org/doc/html/draft-kelly-json-hal-06) specification
- `external_id` (String) The ID of the user in the target app that's linked to the Okta application user object. This value is the native app-specific identifier or primary key for the user in the target app.  The `externa...
- `last_sync` (String) Timestamp of the last synchronization operation. This value is only updated for apps with the `IMPORT_PROFILE_UPDATES` or `PUSH PROFILE_UPDATES` feature.
- `password_changed` (String) Timestamp when the application user password was last changed
- `status` (String) Status of an application user
- `status_changed` (String) Timestamp when the application user status was last changed
- `sync_state` (String) The synchronization state for the application user. The application user's `syncState` depends on whether the `PROFILE_MASTERING` feature is enabled for the app.  > **Note:** User provisioning curr...

## Import

Import using `{app_id}/{id}`:

```shell
terraform import okta_application_user.example <app_id> <id>
```
