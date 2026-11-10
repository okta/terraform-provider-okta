---
page_title: "Data Source: okta_ai_agents_credentials_jwk_ec"
description: |-
  Retrieves an AI agent's elliptic-curve (EC) JSON Web Key by key ID.
---

# Data Source: okta_ai_agents_credentials_jwk_ec

Retrieves an AI agent's elliptic-curve (EC) JSON Web Key by key ID.

## Example Usage

```hcl
data "okta_ai_agents_credentials_jwk_ec" "example" {
  id       = okta_ai_agents_credentials_jwk_ec.example.id
  agent_id = okta_ai_agent.example.id
}
```

## Schema

### Required

- `id` (String) ID of the JSON Web Key.
- `agent_id` (String) ID of the AI agent in Okta.

### Read-Only

- `kid` (String) The unique identifier of the JSON Web Key in the AI agent's JSON Web Key Set (JWKS).
- `alg` (String) Algorithm that's used in the JSON Web Key.
- `crv` (String) The cryptographic curve that's used for the key pair.
- `x` (String) The public x coordinate for the elliptic curve point.
- `y` (String) The public y coordinate for the elliptic curve point.
- `use` (String) Acceptable use of the JSON Web Key.
- `status` (String) Status of the AI agent JSON Web Key.
- `created` (String) Timestamp of when the AI agent JSON Web Key was created.
- `last_updated` (String) Timestamp of when the AI agent JSON Web Key was last updated.
