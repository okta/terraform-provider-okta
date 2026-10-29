---
page_title: "okta_delegation_link Data Source - terraform-provider-okta"
description: |-
  Retrieves a delegation link
---

# okta_delegation_link (Data Source)

Retrieves a delegation link by its unique identifier.

## Example Usage

```terraform
data "okta_delegation_link" "example" {
  id = "dl_abc123def456"
}

output "delegation_link_token_type" {
  value = data.okta_delegation_link.example.from.token_type
}
```

## Schema

### Required

- `id` (String) Unique identifier of the delegation_link.

### Read-Only

- `from` (Block) From (see [below for nested schema](#nestedblock--from))
- `to` (Block) The target resource for the delegation link (see [below for nested schema](#nestedblock--to))

<a id="nestedblock--from"></a>
### Nested Schema for `from`

Read-Only:

- `app_instance_orn` (String) The [ORN](https://developer.
- `client_orn` (String) The [Okta Resource Name (ORN)](/openapi/okta-management/guides/roles/#okta-resource-name-orn) of the OAuth 2.
- `token_type` (String) The type of token accepted by the delegation link.

Optional:

- `type` (String) The type of token source.

<a id="nestedblock--to"></a>
### Nested Schema for `to`

Read-Only:

- `authorization_server_orn` (String) The [Okta Resource Name (ORN)](/openapi/okta-management/guides/roles/#okta-resource-name-orn) of the custom authorization server that protects the target resource.
- `resource_orn` (String) The [Okta Resource Name (ORN)](/openapi/okta-management/guides/roles/#okta-resource-name-orn) of the target resource.