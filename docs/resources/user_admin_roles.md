---
page_title: "Resource: okta_user_admin_roles"
description: |-
  Resource to manage a set of administrator roles for a specific user. This resource allows you to manage admin roles for a single user, independent of the user schema itself.
---

# Resource: okta_user_admin_roles

Resource to manage a set of administrator roles for a specific user. This resource allows you to manage admin roles for a single user, independent of the user schema itself.

## Example Usage

```terraform
resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc-replace_with_uuid@example.com"
  email      = "testAcc-replace_with_uuid@example.com"
}

resource "okta_user_admin_roles" "test" {
  user_id = okta_user.test.id
  admin_roles = [
    "APP_ADMIN",
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `admin_roles` (Set of String) The list of Okta user admin roles, e.g. `['APP_ADMIN', 'USER_ADMIN']` See [API Docs](https://developer.okta.com/docs/api/openapi/okta-management/guides/roles/#standard-roles).
- `user_id` (String) ID of a Okta User

### Optional

- `disable_notifications` (Boolean) When this setting is enabled, the admins won't receive any of the default Okta administrator emails. These admins also won't have access to contact Okta Support and open support cases on behalf of your org.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import okta_user_admin_roles.example <user_id>
```
