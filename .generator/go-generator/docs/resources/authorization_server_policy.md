---
page_title: "okta_authorization_server_policy Resource - terraform-provider-okta"
description: |-
  Manages an Okta Authorization Server Policy resource.
---

# okta_authorization_server_policy (Resource)

Manages an Okta Authorization Server Policy.

## Example Usage

```hcl
resource "okta_authorization_server_policy" "example" {
  auth_server_id = "<auth-server-id>"

  # Optional fields
  # include = "<include>"
  # description = "Example description"
  # name = "Example Name"
  # priority = 0
  # status = "ACTIVE"
}
```

## Schema

### Required

- `auth_server_id` (String) ID of the authorization server

### Optional

- `include` (List) Which clients are included in the policy
- `description` (String) Description of the Policy
- `name` (String) Name of the Policy
- `priority` (Number) Specifies the order in which this Policy is evaluated in relation to the other Policies in a custom authorization server
- `status` (String) Specifies whether requests have access to this Policy
- `system` (Boolean) Specifies whether Okta created this Policy
- `type` (String) Indicates that the Policy is an authorization server Policy

### Read-Only

- `id` (String) The unique identifier for the resource.
- `created` (String) Timestamp when the Policy was created
- `last_updated` (String) Timestamp when the Policy was last updated

## Import

Import using `{auth_server_id}/{id}`:

```shell
terraform import okta_authorization_server_policy.example <auth_server_id> <id>
```
