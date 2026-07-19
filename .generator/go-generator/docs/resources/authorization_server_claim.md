---
page_title: "okta_authorization_server_claim Resource - terraform-provider-okta"
description: |-
  Manages an Okta Authorization Server Claim resource.
---

# okta_authorization_server_claim (Resource)

Manages an Okta Authorization Server Claim.

## Example Usage

```hcl
resource "okta_authorization_server_claim" "example" {
  auth_server_id = "<auth-server-id>"

  # Optional fields
  # always_include_in_token = true
  # claim_type = "<claim_type>"
  # scopes = "<scopes>"
  # group_filter_type = "<group_filter_type>"
  # name = "Example Name"
}
```

## Schema

### Required

- `auth_server_id` (String) ID of the authorization server

### Optional

- `always_include_in_token` (Boolean) Specifies whether to include Claims in the token. The value is always `TRUE` for access token Claims. If the value is set to `FALSE` for an ID token claim, the Claim isn't included in the ID token ...
- `claim_type` (String) Specifies whether the Claim is for an access token (`RESOURCE`) or an ID token (`IDENTITY`)
- `scopes` (List) Scopes
- `group_filter_type` (String) Specifies the type of group filter if `valueType` is `GROUPS`  If `valueType` is `GROUPS`, then the groups returned are filtered according to the value of `group_filter_type`.  If you have complex ...
- `name` (String) Name of the Claim
- `status` (String) Status
- `system` (Boolean) When `true`, indicates that Okta created the Claim
- `value` (String) Specifies the value of the Claim. This value must be a string literal if `valueType` is `GROUPS`, and the string literal is matched with the selected `group_filter_type`. The value must be an Okta ...
- `value_type` (String) Specifies whether the Claim is an Okta Expression Language (EL) expression (`EXPRESSION`), a set of groups (`GROUPS`), or a system claim (`SYSTEM`)

### Read-Only

- `id` (String) The unique identifier for the resource.

## Import

Import using `{auth_server_id}/{id}`:

```shell
terraform import okta_authorization_server_claim.example <auth_server_id> <id>
```
