---
page_title: "Resource: okta_a2a_servers_authorization_server"
description: |-
  Links an authorization server to an Okta A2A (Agent-to-Agent) server.
---

# Resource: okta_a2a_servers_authorization_server

Links an authorization server to an Okta Agent-to-Agent (A2A) server. The A2A Server Registrations V2 API provides read-only operations for A2A servers and management operations for their associated authorization servers.

Updates are not supported in place. Changing `a2a_server_id` or `orn` replaces the link.

## Example Usage

```hcl
resource "okta_a2a_servers_authorization_server" "example" {
  a2a_server_id = "<a2a-server-id>"
  orn           = "<authorization-server-orn>"
}
```

## Schema

### Required

- `a2a_server_id` (String) ID of an A2A server. Forces replacement if changed.
- `orn` (String) The [ORN](https://developer.okta.com/docs/reference/okta-resource-names/) of the authorization server to link.

### Read-Only

- `id` (String) The unique identifier of the link, which is the authorization server's ID.
- `authorization_server` (Object) An authorization server. Attributes:
  - `id` (String) Unique identifier for the authorization server.
  - `issuer` (String) OAuth 2.0 issuer of the authorization server.
  - `orn` (String) The [ORN](https://developer.okta.com/docs/reference/okta-resource-names/) of the authorization server.
  - `status` (String) Current status of the authorization server in its lifecycle.
  - `subject_mapping_expression` (String) Okta Expression Language formula that computes the `sub` claim value for the ID-JAG token minted for this authorization server.
  - `type` (String) How the authorization server's metadata is sourced.

## Import

Import using `{a2a_server_id}/{id}`:

```shell
terraform import okta_a2a_servers_authorization_server.example <a2a-server-id>/<authorization-server-id>
```

> **Note**: `orn` is a configuration input and is not returned by the read endpoint, so it is not restored on import. Set it in configuration after importing.
