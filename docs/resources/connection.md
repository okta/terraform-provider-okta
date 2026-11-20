---
page_title: "Resource: okta_connection"
description: |-
  Manages a resource connection between a workload principal and a resource it can access.
---

# Resource: okta_connection

Manages a resource connection. A connection lets a workload principal (the actor, such as an AI agent) obtain a credential for a target resource, such as an MCP server, API server, or A2A server. The Resource Connections V2 API manages these connections.

The provider activates the connection after creating it, and deactivates it before deleting it. Changes to `protocol.scopes`, `protocol.scope_condition`, and `resource_identifier` are applied in place.

## Example Usage

An AI agent connecting to an MCP server through an Okta authorization server, using identity assertion:

```hcl
resource "okta_ai_agent" "example" {
  profile {
    name        = "Example agent"
    description = "Agent that calls the Teams MCP server"
  }

  resource_url = "https://agent.example.com/api"
}

resource "okta_connection" "example" {
  actor {
    wlp_id = okta_ai_agent.example.id
  }

  authorization_server {
    orn = "orn:okta:idp:<org-id>:authorization_servers:<authorization-server-id>"
  }

  protocol {
    protocol_type   = "IDENTITY_ASSERTION"
    scope_condition = "INCLUDE_ONLY"
    scopes          = ["openid", "profile", "email"]
  }

  resource {
    resource_type = "MCP_SERVER"
    orn           = "orn:okta:directory:<org-id>:resource-servers:mcp:<mcp-server-id>"
  }
}
```

The authorization server must already be linked to the target resource, for example with [`okta_mcp_servers_authorization_server`](mcp_servers_authorization_server.md) for an MCP server.

## Schema

### Optional

- `resource_identifier` (String) The resource identifier (audience) used when requesting tokens for the resource. If unset, Okta fills it in.
- `actor` (Block) The workload principal that uses the connection to reach the resource.
  - `wlp_id` (String, Required) ID of the workload principal, for example an AI agent or an MCP server.
  - `display_name` (String, Read-Only) Display name of the workload principal.
  - `wlp_type` (String, Read-Only) Type of the workload principal, for example `AI_AGENT`.
- `authorization_server` (Block) The authorization server where the actor redeems the identity assertion.
  - `orn` (String) The [ORN](https://developer.okta.com/docs/api/openapi/okta-management/guides/roles/#okta-resource-name-orn) of the authorization server.
  - `display_name` (String, Read-Only) Display name of the authorization server.
  - `issuer_url` (String, Read-Only) Issuer URL of the authorization server.
- `protocol` (Block) How the actor obtains the credential that it presents to the resource.
  - `protocol_type` (String, Required) How the actor obtains the credential: `IDENTITY_ASSERTION` or `STS`.
  - `scope_condition` (String) How Okta evaluates requested scopes for the connection, for example `INCLUDE_ONLY` or `ALL_SCOPES`.
  - `scopes` (List of String) The scopes the actor can request.
  - `client_id_at_resource` (String) The OAuth client identifier that the actor is registered as at the external authorization server.
  - `credential` (Block) The credential that the STS uses to obtain a token for the resource. Used with the `STS` protocol type.
    - `credential_type` (String) Type of the credential. Defaults to `ACCESS_TOKEN`.
    - `orn` (String) The ORN of the credential.
    - `display_name` (String, Read-Only) Display name of the credential.
    - `scopes` (List of String, Read-Only) The OAuth scopes configured on the client auth settings that back this credential.
- `resource` (Block) The target resource of the connection, identified by its ORN.
  - `resource_type` (String) Type of the target resource: `MCP_SERVER`, `API_SERVER`, `A2A_SERVER`, `SERVICE_ACCOUNT`, or `VAULT_SECRET`.
  - `orn` (String) The ORN of the target resource.
  - `display_name` (String, Read-Only) Display name of the resource.
  - `description` (String, Read-Only) Description of the resource. Set for vaulted secrets.
  - `location` (String, Read-Only) Location of the resource, for example the MCP server URL or the secret's path in Okta Privileged Access.
  - `app_instance` (Block) The app instance that owns the resource.
    - `orn` (String) The ORN of the app instance.
    - `display_name` (String, Read-Only) Display name of the app instance.
    - `logo` (String, Read-Only) Image URL for the app logo.

### Read-Only

- `id` (String) The unique identifier of the connection.
- `orn` (String) The ORN of the connection.
- `status` (String) Status of the connection, for example `ACTIVE`.

## Import

Import using the connection ID:

```shell
terraform import okta_connection.example <connection-id>
```
