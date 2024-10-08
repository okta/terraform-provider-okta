---
page_title: "Resource: okta_group_owner"
description: |-
  Manages group owner resource.
---

# Resource: okta_group_owner

Manages group owner resource.

## Example Usage

```terraform
resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc-replace_with_uuid@example.com"
  email      = "testAcc-replace_with_uuid@example.com"
}

resource "okta_group" "test" {
  name = "testAcc_replace_with_uuid"
}

resource "okta_group_owner" "test" {
  group_id          = okta_group.test.id
  id_of_group_owner = okta_user.test.id
  type              = "USER"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `group_id` (String) The id of the group
- `id_of_group_owner` (String) The user id of the group owner
- `type` (String) The entity type of the owner. Enum: "GROUP" "USER"

### Read-Only

- `display_name` (String) The display name of the group owner
- `id` (String) The id of the group owner resource
- `origin_id` (String) The ID of the app instance if the originType is APPLICATION. This value is NULL if originType is OKTA_DIRECTORY.
- `origin_type` (String) The source where group ownership is managed. Enum: "APPLICATION" "OKTA_DIRECTORY"
- `resolved` (Boolean) If originType is APPLICATION, this parameter is set to FALSE until the owner's originId is reconciled with an associated Okta ID.


