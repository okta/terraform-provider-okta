---
page_title: "okta_realm Resource - terraform-provider-okta"
description: |-
  Manages an Okta Realm resource.
---

# okta_realm (Resource)

Manages an Okta Realm.

## Example Usage

```hcl
resource "okta_realm" "example" {
  name = "Example Name"

  # Optional fields
  # domains = "<domains>"
  # realm_type = "<realm_type>"
}
```

## Schema

### Required

- `name` (String) Name of a realm

### Optional

- `domains` (List) Array of allowed domains. No user in this realm can be created or updated unless they have a username and email from one of these domains.  The following characters aren't allowed in the domain nam...
- `realm_type` (String) Used to store partner users. This property must be set to `PARTNER` to access Okta's external partner portal.

### Read-Only

- `id` (String) The unique identifier for the resource.
- `created` (String) Timestamp when the realm was created
- `is_default` (Boolean) Indicates the default realm. Existing users will start out in the default realm and can be moved to other realms individually or through realm assignments. See [Realms Assignments API](/openapi/okt...
- `last_updated` (String) Timestamp when the realm was updated
