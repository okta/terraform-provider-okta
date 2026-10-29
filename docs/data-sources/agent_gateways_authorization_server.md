---
page_title: "Data Source: okta_agent_gateways_authorization_server"
description: |-
  Retrieves an authorization server linked to an agent gateway
---

# Data Source: okta_agent_gateways_authorization_server

Retrieves a specific authorization server that is linked to an agent gateway by ID. Authorization servers define the OAuth 2.0 endpoints and credentials needed for agent authentication.

## Example Usage

```terraform
data "okta_agent_gateways_authorization_server" "example" {
  id               = "as_abc123def456"
  agent_gateway_id = "agw_xyz789"
}

output "authorization_server" {
  value = {
    id     = data.okta_agent_gateways_authorization_server.example.id
    issuer = data.okta_agent_gateways_authorization_server.example.issuer
    status = data.okta_agent_gateways_authorization_server.example.status
  }
}
```

## Argument Reference

### Required

- `id` - (String) Unique identifier of the authorization server
- `agent_gateway_id` - (String) ID of the agent gateway

## Attribute Reference

- `issuer` - (String) OAuth 2.0 issuer URL for the authorization server
- `last_updated` - (String) Timestamp when the authorization server was last updated (RFC3339 format)
- `orn` - (String) The [ORN](https://developer.okta.com/docs/concepts/orn/) identifier
- `status` - (String) Current status of the authorization server
- `metadata` - (Block) OAuth 2.0 metadata about the authorization server
  - `authorization_endpoint` - (String) URL of the authorization endpoint
  - `grant_types_supported` - (List of String) Supported OAuth 2.0 grant types
  - `token_endpoint` - (String) URL of the token endpoint
  - `token_endpoint_auth_method` - (String) Client authentication method for the token endpoint

## Example: Retrieve Authorization Server Metadata

```terraform
data "okta_agent_gateways_authorization_server" "oauth" {
  id               = "as_oauth_server"
  agent_gateway_id = "agw_prod"
}

output "oauth_endpoints" {
  value = {
    issuer               = data.okta_agent_gateways_authorization_server.oauth.issuer
    authorization_url    = try(data.okta_agent_gateways_authorization_server.oauth.metadata[0].authorization_endpoint, "")
    token_url            = try(data.okta_agent_gateways_authorization_server.oauth.metadata[0].token_endpoint, "")
    auth_method          = try(data.okta_agent_gateways_authorization_server.oauth.metadata[0].token_endpoint_auth_method, "")
    supported_grant_types = try(data.okta_agent_gateways_authorization_server.oauth.metadata[0].grant_types_supported, [])
  }
}
```

## Example: Check Authorization Server Status

```terraform
variable "gateway_id" {
  type = string
}

variable "auth_server_id" {
  type = string
}

data "okta_agent_gateways_authorization_server" "check" {
  id               = var.auth_server_id
  agent_gateway_id = var.gateway_id
}

locals {
  is_active = data.okta_agent_gateways_authorization_server.check.status == "ACTIVE"
}

output "auth_server_status" {
  value = {
    id            = data.okta_agent_gateways_authorization_server.check.id
    active        = local.is_active
    last_updated  = data.okta_agent_gateways_authorization_server.check.last_updated
  }
}
```

## Notes

- This data source retrieves authorization servers linked to a specific agent gateway
- Authorization servers must be in an `ACTIVE` state to authenticate agents
- The metadata block contains OAuth 2.0 configuration details needed for credential exchanges
- Use the ORN field to reference the authorization server in other Okta resources