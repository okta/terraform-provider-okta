---
layout: 'okta' page_title: 'Okta: okta_auth_server_claim_default' sidebar_current: '
docs-okta-resource-auth-server-claim-default' description: |- Configures Default Authorization Server Claim
---

# okta_auth_server_claim_default

Configures Default Authorization Server Claim.

This resource allows you to configure Default Authorization Server Claims.

## Example Usage

```hcl
resource "okta_auth_server_claim_default" "example" {
  auth_server_id = "<auth server id>"
  name           = "sub"
  value          = "(appuser != null) ? appuser.userName : app.clientId"
}
```

## Argument Reference

The following arguments are supported:

- `auth_server_id` - (Required) ID of the authorization server.

- `name` - (Required) The name of the claim. Can be set to `"sub"`, `"address"`, `"birthdate"`, `"email"`,
  `"email_verified"`, `"family_name"`, `"gender"`, `"given_name"`, `"locale"`, `"middle_name"`, `"name"`, `"nickname"`,
  `"phone_number"`, `"picture"`, `"preferred_username"`, `"profile"`, `"updated_at"`, `"website"`, `"zoneinfo"`.
  
- `value` - (Optional/Required) The value of the claim. Only required for `"sub"` claim.

## Attributes Reference

- `id` - The ID for the auth server claim.

- `name` - The name of the claim.

- `scopes` - The list of scopes the auth server claim is tied to.

- `status` - The status of the application.

- `value` - The value of the claim.

- `value_type` - The type of value of the claim.

- `claim_type` - Specifies whether the claim is for an access token `"RESOURCE"` or ID token `"IDENTITY"`.

- `always_include_in_token` - Specifies whether to include claims in token.

## Import

Authorization Server Claim can be imported via the Auth Server ID and Claim ID or Claim Name.

```
$ terraform import okta_auth_server_claim_default.example <auth server id>/<claim id>
```

or

```
$ terraform import okta_auth_server_claim_default.example <auth server id>/<claim name>
```
