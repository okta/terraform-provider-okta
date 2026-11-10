---
page_title: "Resource: okta_ai_agents_credentials_jwk_ec"
description: |-
  Adds an elliptic-curve (EC) JSON Web Key to an AI agent.
---

# Resource: okta_ai_agents_credentials_jwk_ec

Adds an elliptic-curve (EC) JSON Web Key to an AI agent. Okta uses the key to verify tokens the agent signs.

On destroy the key is deactivated and then deleted.

~> **Note:** A key can't be changed after it is created, and that includes `status`. Changing `agent_id` or `kty` replaces the key; changing any other attribute isn't supported and fails the apply with "Provider produced inconsistent result". To use different key material, add a new key resource.

## Example Usage

```hcl
resource "okta_ai_agent" "example" {
  profile {
    name = "example-agent"
  }

  resource_url = "https://example.com/api"

  sign_on_provider {
    type = "NEW_OIDC_APP"
  }
}

resource "okta_ai_agents_credentials_jwk_ec" "example" {
  agent_id = okta_ai_agent.example.id

  kty    = "EC"
  alg    = "ES256"
  kid    = "agent-signing-key"
  crv    = "P-256"
  use    = "sig"
  status = "ACTIVE"

  x = "WKn-ZIGevcwGIyyrzFoZNBdaq9_TsqzGl96oc0CWuis"
  y = "y77t-RvAHRKTsSGdIYUfweuOvwrvDD-Q3Hv5J0fSKQE"
}
```

## Schema

### Required

- `agent_id` (String) ID of the AI agent in Okta. Forces replacement if changed.
- `kty` (String) Key type. Must be `EC`. Forces replacement if changed.

### Optional

- `kid` (String) The unique identifier of the JSON Web Key in the AI agent's JSON Web Key Set (JWKS).
- `alg` (String) Algorithm that's used in the JSON Web Key.
- `use` (String) Acceptable use of the JSON Web Key.
- `status` (String) Status of the AI agent JSON Web Key: `ACTIVE` or `INACTIVE`. Set at creation; it can't be changed afterwards.
- `crv` (String) The cryptographic curve that's used for the key pair, for example `P-256`.
- `x` (String) The public x coordinate for the elliptic curve point, Base64URL-encoded.
- `y` (String) The public y coordinate for the elliptic curve point, Base64URL-encoded.

### Read-Only

- `id` (String) ID of the JSON Web Key.

## Import

Import using `{agent_id}/{id}`:

```shell
terraform import okta_ai_agents_credentials_jwk_ec.example <agent-id>/<key-id>
```
