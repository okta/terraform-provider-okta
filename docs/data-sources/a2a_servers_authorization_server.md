---
page_title: "Data Source: okta_a2a_servers_authorization_server"
description: |-
  Retrieves the link between an A2A server and an authorization server, by the authorization server's ID.
---

# Data Source: okta_a2a_servers_authorization_server

Retrieves the link between an Okta Agent-to-Agent (A2A) server and an authorization server, by the authorization server's ID.

## Example Usage

```hcl
data "okta_a2a_servers_authorization_server" "example" {
  id            = "<authorization-server-id>"
  a2a_server_id = "<a2a-server-id>"
}
```

## Schema

### Required

- `id` (String) Unique identifier of the authorization server linked to the A2A server.
- `a2a_server_id` (String) ID of an A2A server.

### Read-Only

- `created` (String) Timestamp when the authorization server was linked to the resource server.
- `last_updated` (String) Timestamp when the link was last updated.
- `authorization_server` (Block) An authorization server. See [below](#nested-schema-for-authorization_server).

<a id="nested-schema-for-authorization_server"></a>
### Nested Schema for `authorization_server`

Read-Only:

- `created` (String) Timestamp when the authorization server was created.
- `id` (String) Unique identifier for the authorization server.
- `issuer` (String) OAuth 2.0 issuer of the authorization server.
- `last_updated` (String) Timestamp when the authorization server was last updated.
- `orn` (String) The [ORN](https://developer.okta.com/docs/reference/okta-resource-names/) of the authorization server.
- `status` (String) Current status of the authorization server in its lifecycle.
- `subject_mapping_expression` (String) Okta Expression Language formula that computes the `sub` claim value for the ID-JAG token minted for this authorization server.
- `type` (String) How the authorization server's metadata is sourced.
- `metadata` (Block) Metadata about the authorization server. See [below](#nested-schema-for-authorization_servermetadata).

<a id="nested-schema-for-authorization_servermetadata"></a>
### Nested Schema for `authorization_server.metadata`

Read-Only:

- `authorization_endpoint` (String) URL of the authorization endpoint.
- `grant_types_supported` (List of String) OAuth 2.0 grant types supported by the authorization server.
- `identity_chaining_requested_token_types_supported` (List of String) Token types supported for identity chaining requested tokens.
- `protocol_supported` (List of String) The mechanism(s) this authorization server supports for obtaining tokens for resource servers, some combination of `STS` and `IDENTITY_ASSERTION`.
- `registration_endpoint` (String) URL of the OAuth 2.0 dynamic client registration endpoint.
- `token_endpoint` (String) URL of the token endpoint.
- `token_endpoint_auth_method` (String) The client authentication method used when presenting tokens to the authorization server's token endpoint.
- `token_endpoint_auth_methods_supported` (List of String) Client authentication methods supported by the token endpoint.
