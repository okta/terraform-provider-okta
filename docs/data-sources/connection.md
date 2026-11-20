---
page_title: "Data Source: okta_connection"
description: |-
  Retrieves a resource connection between a workload principal and a resource it can access.
---

# Data Source: okta_connection

Retrieves a resource connection by ID. A connection lets a workload principal (the actor, such as an AI agent) obtain a credential for a target resource, such as an MCP server, API server, or A2A server.

## Example Usage

```terraform
data "okta_connection" "example" {
  id = "<connection-id>"
}

output "connection_target" {
  value = {
    actor    = data.okta_connection.example.actor.wlp_id
    protocol = data.okta_connection.example.protocol.protocol_type
    resource = data.okta_connection.example.resource.display_name
  }
}
```

## Argument Reference

### Required

- `id` (String) Unique identifier of the connection.

## Attributes Reference

- `orn` (String) The [ORN](https://developer.okta.com/docs/api/openapi/okta-management/guides/roles/#okta-resource-name-orn) of the connection.
- `status` (String) Status of the connection, for example `ACTIVE`.
- `resource_identifier` (String) The resource identifier (audience) used when requesting tokens for the resource.
- `actor` (Block) The workload principal that uses the connection to reach the resource.
  - `wlp_id` (String) ID of the workload principal, for example an AI agent or an MCP server.
  - `wlp_type` (String) Type of the workload principal, for example `AI_AGENT`.
  - `display_name` (String) Display name of the workload principal.
- `authorization_server` (Block) The authorization server where the actor redeems the identity assertion.
  - `orn` (String) The ORN of the authorization server.
  - `display_name` (String) Display name of the authorization server.
  - `issuer_url` (String) Issuer URL of the authorization server.
- `protocol` (Block) How the actor obtains the credential that it presents to the resource.
  - `protocol_type` (String) How the actor obtains the credential: `IDENTITY_ASSERTION` or `STS`.
  - `scope_condition` (String) How Okta evaluates requested scopes for the connection, for example `INCLUDE_ONLY` or `ALL_SCOPES`.
  - `scopes` (List of String) The scopes the actor can request.
  - `client_id_at_resource` (String) The OAuth client identifier that the actor is registered as at the external authorization server.
  - `credential` (Block) The credential that the STS uses to obtain a token for the resource. Set for the `STS` protocol type.
    - `credential_type` (String) Type of the credential, for example `ACCESS_TOKEN`.
    - `orn` (String) The ORN of the credential.
    - `display_name` (String) Display name of the credential.
    - `scopes` (List of String) The OAuth scopes configured on the client auth settings that back this credential.
- `resource` (Block) The target resource of the connection.
  - `resource_type` (String) Type of the target resource: `MCP_SERVER`, `API_SERVER`, `A2A_SERVER`, `SERVICE_ACCOUNT`, or `VAULT_SECRET`.
  - `orn` (String) The ORN of the target resource.
  - `display_name` (String) Display name of the resource.
  - `description` (String) Description of the resource. Set for vaulted secrets.
  - `location` (String) Location of the resource, for example the MCP server URL or the secret's path in Okta Privileged Access.
  - `app_instance` (Block) The app instance that owns the resource.
    - `orn` (String) The ORN of the app instance.
    - `display_name` (String) Display name of the app instance.
    - `logo` (String) Image URL for the app logo.
