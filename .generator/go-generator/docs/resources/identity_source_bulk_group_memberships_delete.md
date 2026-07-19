---
page_title: "okta_identity_source_bulk_group_memberships_delete Resource - terraform-provider-okta"
description: |-
  Stages group memberships to be removed from groups in Okta within an identity source session.
---

# okta_identity_source_bulk_group_memberships_delete (Resource)

Stages group memberships to be removed from groups in Okta as part of an identity source session
(`POST /api/v1/identity-sources/{identitySourceId}/sessions/{sessionId}/bulk-group-memberships-delete`).

This resource corresponds to a single API call that queues removals into the session.
It can be used alongside `okta_identity_source_bulk_group_memberships_upsert` in the same session
to stage both additions and removals before the session is committed.

> **Note:** This is a write-only staging operation. There is no read or delete endpoint for staged
> data. Destroying this Terraform resource does not remove data from the session. The `id` is set
> to the `session_id` value after a successful call.

## Example Usage

```hcl
resource "okta_identity_source_session" "example" {
  identity_source_id = "<identity-source-id>"
}

resource "okta_identity_source_bulk_group_memberships_delete" "example" {
  identity_source_id = okta_identity_source_session.example.identity_source_id
  session_id         = okta_identity_source_session.example.id

  memberships = [
    {
      group_external_id   = "group-ext-id-1"
      member_external_ids = ["user-ext-id-1", "user-ext-id-2"]
    },
    {
      group_external_id   = "group-ext-id-2"
      member_external_ids = ["user-ext-id-3"]
    }
  ]
}
```

### Staging both additions and removals in the same session

```hcl
resource "okta_identity_source_session" "example" {
  identity_source_id = "<identity-source-id>"
}

resource "okta_identity_source_bulk_group_memberships_upsert" "add" {
  identity_source_id = okta_identity_source_session.example.identity_source_id
  session_id         = okta_identity_source_session.example.id

  memberships = [
    {
      group_external_id   = "group-ext-id-1"
      member_external_ids = ["user-ext-id-new"]
    }
  ]
}

resource "okta_identity_source_bulk_group_memberships_delete" "remove" {
  identity_source_id = okta_identity_source_session.example.identity_source_id
  session_id         = okta_identity_source_session.example.id

  memberships = [
    {
      group_external_id   = "group-ext-id-1"
      member_external_ids = ["user-ext-id-old"]
    }
  ]
}
```

## Schema

### Required

- `identity_source_id` (String) ID of the identity source.
- `session_id` (String) ID of the identity source session. Forces replacement when changed.

### Optional

- `memberships` (List of Object) Array of group memberships to stage for removal in Okta.
  Up to 200 items per request.
  - `group_external_id` (String) The external ID of the group whose memberships need to be deleted
    in Okta.
  - `member_external_ids` (List of String) External IDs of member profiles to remove from this group
    in Okta.

### Read-Only

- `id` (String) Set to the `session_id` value after a successful staging call.

## Import

Import using `{identity_source_id}/{session_id}/{id}`:

```shell
terraform import okta_identity_source_bulk_group_memberships_delete.example "<identity_source_id>/<session_id>/<id>"
```
