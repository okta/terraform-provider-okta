# okta_group_owner

This resource represents a Group Owner for an Okta organization. More information can
be found in the
[Group Owners](https://developer.okta.com/docs/api/openapi/okta-management/management/tag/GroupOwner/#tag/GroupOwner) API
documentation.

- Example [resource.tf](./resource.tf)

## Import

An okta_group_owner resource can be imported using the following format:

```bash
terraform import okta_group_owner.example group_id/group_owner_id
```

Where:

- `group_id` is the ID of the group
- `group_owner_id` is the ID of the group owner resource

**Example:**

```bash
terraform import okta_group_owner.example group_123/group_owner_456
```

**Note:** When importing, you must still provide the required `group_id`, `id_of_group_owner`, and `type` attributes in your Terraform configuration, as these are not stored in the import ID.
