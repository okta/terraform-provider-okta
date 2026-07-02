---
page_title: "Resource: okta_identity_source_import"
description: |-
  Runs a complete identity source import job in a single resource.
---

# Resource: okta_identity_source_import

Orchestrates a full Okta Identity Source import job in one resource. On each apply it:

1. Creates a new import session.
2. Uploads any combination of staged data (upsert/delete for users, groups, and group memberships).
3. Triggers the import via `startImportFromIdentitySource`.

If any upload step fails the session is automatically deleted so a corrected re-apply is not blocked by Okta's rate limit (*"Only one active import session per identity source is allowed to be created every 5 minutes"*).

~> **Note:** This resource does not support deletion. Removing it from configuration will emit a warning but does not undo the import in Okta. To re-run the import, change at least one attribute so Terraform creates a new resource.

## Example Usage

### Upsert users and group memberships

```terraform
resource "okta_identity_source_import" "example" {
  identity_source_id = "<identity-source-id>"

  upsert_users {
    entity_type = "USERS"

    profiles {
      external_id = "USEREXT001"

      profile {
        user_name  = "jdoe@example.com"
        email      = "jdoe@example.com"
        first_name = "Jane"
        last_name  = "Doe"
      }
    }
  }

  upsert_group_memberships {
    memberships {
      group_external_id   = "GROUPEXT001"
      member_external_ids = ["USEREXT001"]
    }
  }
}
```

### Upsert and delete in the same job

```terraform
resource "okta_identity_source_import" "example" {
  identity_source_id = "<identity-source-id>"

  upsert_users {
    entity_type = "USERS"

    profiles {
      external_id = "USEREXT_NEW"

      profile {
        user_name  = "new.user@example.com"
        email      = "new.user@example.com"
        first_name = "New"
        last_name  = "User"
      }
    }
  }

  delete_users {
    entity_type = "USERS"

    profiles {
      external_id = "USEREXT_OLD"
    }
  }

  upsert_groups {
    profiles {
      external_id = "GROUPEXT001"

      group_profile {
        display_name = "Engineering"
        description  = "Engineering team"
      }
    }
  }

  delete_groups {
    external_ids = ["GROUPEXT_DEPRECATED"]
  }
}
```

## Schema

### Required

- `identity_source_id` (String) ID of the identity source. Forces replacement when changed.

### Read-Only

- `id` (String) Session ID of the triggered import job.
- `session_id` (String) Session ID created for this import job.
- `session_status` (String) Status of the import session after triggering (e.g. `IN_PROGRESS`, `COMPLETED`).

### Optional Blocks

All upload blocks are optional. At least one should be provided.

- `upsert_users` (Block, Optional) Users to create or update. (see [below](#nested-schema-for-upsert_users))
- `upsert_groups` (Block, Optional) Groups to create or update. (see [below](#nested-schema-for-upsert_groups))
- `delete_users` (Block, Optional) Users to delete from Okta. (see [below](#nested-schema-for-delete_users))
- `delete_groups` (Block, Optional) Groups to delete from Okta. (see [below](#nested-schema-for-delete_groups))
- `upsert_group_memberships` (Block, Optional) Group memberships to create or update. (see [below](#nested-schema-for-upsert_group_memberships))
- `delete_group_memberships` (Block, Optional) Group memberships to delete. (see [below](#nested-schema-for-delete_group_memberships))

---

### Nested Schema for `upsert_users`

Optional:

- `entity_type` (String) Entity type. Currently only `USERS` is supported.
- `profiles` (Block List) User profiles to upsert. (see [below](#nested-schema-for-upsert_usersprofiles))

### Nested Schema for `upsert_users.profiles`

Optional:

- `external_id` (String) External ID of the user.
- `profile` (Block, Optional) User profile attributes. (see [below](#nested-schema-for-upsert_usersprofilesprofile))

### Nested Schema for `upsert_users.profiles.profile`

Optional:

- `email` (String) Email address of the user.
- `first_name` (String) First name of the user.
- `home_address` (String) Home address of the user.
- `last_name` (String) Last name of the user.
- `mobile_phone` (String) Mobile phone number of the user.
- `second_email` (String) Alternative email address of the user.
- `user_name` (String) Username of the user.

---

### Nested Schema for `upsert_groups`

Optional:

- `profiles` (Block List) Group profiles to upsert. (see [below](#nested-schema-for-upsert_groupsprofiles))

### Nested Schema for `upsert_groups.profiles`

Optional:

- `external_id` (String) External ID of the group.
- `group_profile` (Block, Optional) Group profile attributes. (see [below](#nested-schema-for-upsert_groupsprofilesgroup_profile))

### Nested Schema for `upsert_groups.profiles.group_profile`

Optional:

- `description` (String) Description of the group.
- `display_name` (String) Display name of the group.

---

### Nested Schema for `delete_users`

Optional:

- `entity_type` (String) Entity type. Currently only `USERS` is supported.
- `profiles` (Block List) User profiles to delete (by external ID). (see [below](#nested-schema-for-delete_usersprofiles))

### Nested Schema for `delete_users.profiles`

Optional:

- `external_id` (String) External ID of the user to delete.

---

### Nested Schema for `delete_groups`

Optional:

- `external_ids` (List of String) External IDs of the groups to delete.

---

### Nested Schema for `upsert_group_memberships`

Optional:

- `memberships` (Block List) Group memberships to upsert. (see [below](#nested-schema-for-upsert_group_membershipsmemberships))

### Nested Schema for `upsert_group_memberships.memberships`

Optional:

- `group_external_id` (String) External ID of the group.
- `member_external_ids` (List of String) External IDs of the group members to add.

---

### Nested Schema for `delete_group_memberships`

Optional:

- `memberships` (Block List) Group memberships to delete. (see [below](#nested-schema-for-delete_group_membershipsmemberships))

### Nested Schema for `delete_group_memberships.memberships`

Optional:

- `group_external_id` (String) External ID of the group.
- `member_external_ids` (List of String) External IDs of the group members to remove.