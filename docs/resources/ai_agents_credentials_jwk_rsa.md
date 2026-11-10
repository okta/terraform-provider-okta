---
page_title: "Resource: okta_ai_agents_credentials_jwk_rsa"
description: |-
  Adds an RSA JSON Web Key to an AI agent.
---

# Resource: okta_ai_agents_credentials_jwk_rsa

Adds an RSA JSON Web Key to an AI agent. Okta uses the key to verify tokens the agent signs.

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

resource "okta_ai_agents_credentials_jwk_rsa" "example" {
  agent_id = okta_ai_agent.example.id

  kty    = "RSA"
  alg    = "RS256"
  kid    = "agent-signing-key"
  use    = "sig"
  status = "ACTIVE"

  n = "<base64url-encoded modulus>"
  e = "AQAB"
}
```

## Schema

### Required

- `agent_id` (String) ID of the AI agent in Okta. Forces replacement if changed.
- `kty` (String) Key type. Must be `RSA`. Forces replacement if changed.

### Optional

- `kid` (String) The unique identifier of the JSON Web Key in the AI agent's JSON Web Key Set (JWKS).
- `alg` (String) Algorithm that's used in the JSON Web Key.
- `use` (String) Acceptable use of the JSON Web Key.
- `status` (String) Status of the AI agent JSON Web Key: `ACTIVE` or `INACTIVE`. Set at creation; it can't be changed afterwards.
- `n` (String) The modulus of the RSA public key, represented as a Base64URL-encoded string.
- `e` (String) The public exponent of the RSA key, represented as a Base64URL-encoded string, for example `AQAB`.

### Read-Only

- `id` (String) ID of the JSON Web Key.

## Import

Import using `{agent_id}/{id}`:

```shell
terraform import okta_ai_agents_credentials_jwk_rsa.example <agent-id>/<key-id>
```
