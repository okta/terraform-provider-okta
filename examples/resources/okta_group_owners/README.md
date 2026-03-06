# okta_group_owners

Manage all or a tracked subset of owners for an Okta group.

```hcl
resource "okta_group_owners" "owners" {
  group_id = okta_group.grp.id

  owner {
    type = "USER"
    id   = okta_user.owner1.id
  }
  owner {
    type = "USER"
    id   = okta_user.owner2.id
  }

  # Optional
  # track_all_owners = false
}
```
