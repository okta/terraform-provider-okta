# okta_group_owners

Manage all owners for an Okta group. The resource is authoritative: any owners not declared in configuration will be removed.

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
}
```
