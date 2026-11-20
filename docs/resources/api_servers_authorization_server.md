---
page_title: "Resource: okta_api_servers_authorization_server"
description: |-
  Links an authorization server to an Okta API server.
---

# Resource: okta_api_servers_authorization_server

Links an authorization server to an Okta API server. The API Server Registrations V2 API registers third-party API servers as resource servers, and this resource manages the authorization servers linked to them.

The link has no updatable fields. Changing `api_server_id` replaces the link.

## Example Usage

```hcl
resource "okta_api_server" "example" {
  resource_url = "https://api.example.com"

  metadata {
    display_name     = "Example API"
    scopes_supported = ["read:items"]
  }
}

resource "okta_api_servers_authorization_server" "example" {
  api_server_id = okta_api_server.example.id
  orn           = "orn:okta:directory:<org-id>:external-authorization-servers:<authorization-server-id>"
}
```

## Schema

### Required

- `api_server_id` (String) ID of an API server. Forces replacement if changed.
- `orn` (String) The [ORN](https://developer.okta.com/docs/api/openapi/okta-management/guides/roles/#okta-resource-name-orn) of an existing Okta custom or external authorization server, for example `orn:okta:directory:<org-id>:external-authorization-servers:<authorization-server-id>`.

### Read-Only

- `id` (String) The unique identifier of the link, which is the authorization server's ID.
- `tenant_id` (String) Identifier of the tenant at the authorization server for the app instance, emitted as the `aud_tenant` claim in the ID-JAG token.
- `authorization_server` (Object) The linked authorization server. Attributes:
  - `id` (String) Unique identifier for the authorization server.
  - `issuer` (String) OAuth 2.0 issuer of the authorization server.
  - `orn` (String) The ORN of the authorization server.
  - `status` (String) Current status of the authorization server in its lifecycle.
  - `subject_mapping_expression` (String) Okta Expression Language formula that computes the `sub` claim value for the ID-JAG token minted for this authorization server.
  - `type` (String) How the authorization server's metadata is sourced.

## Import

Import using `{api_server_id}/{id}`:

```shell
terraform import okta_api_servers_authorization_server.example <api-server-id>/<authorization-server-id>
```

> **Note**: `orn` is a configuration input and is not returned by the read endpoint, so it is not restored on import. Set it in configuration after importing.
