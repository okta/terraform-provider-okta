---
page_title: "okta_realm_assignment Resource - terraform-provider-okta"
description: |-
  Manages an Okta Realm Assignment resource.
---

# okta_realm_assignment (Resource)

Manages an Okta Realm Assignment.

## Example Usage

```hcl
resource "okta_realm_assignment" "example" {

  # Optional fields
  # realm_id = "<realm-id>"
  # value = "<value>"
  # profile_source_id = "<profile-source-id>"
  # domains = "<domains>"
  # name = "Example Name"
}
```

## Schema

### Optional

- `realm_id` (String) ID of the realm
- `value` (String) Value of the condition expression
- `profile_source_id` (String) ID of the profile source
- `domains` (List) Array of allowed domains. No user in this realm can be created or updated unless they have a username and email from one of these domains.  The following characters aren't allowed in the domain nam...
- `name` (String) Name of the realm
- `priority` (Number) The priority of the realm assignment. The lower the number, the higher the priority. This helps resolve conflicts between realm assignments. > **Note:** When you create realm assignments in bulk, r...
- `status` (String) Status

### Read-Only

- `id` (String) The unique identifier for the resource.
- `created` (String) Timestamp when the realm assignment was created
- `is_default` (Boolean) Indicates the default realm. Existing users will start out in the default realm and can be moved individually to other realms.
- `last_updated` (String) Timestamp of when the realm assignment was updated
