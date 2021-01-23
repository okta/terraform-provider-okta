---
layout: 'okta'
page_title: 'Okta: okta_idp_oidc'
sidebar_current: 'docs-okta-datasource-idp-oidc'
description: |-
  Get a OIDC IdP from Okta.
---

# okta_idp_oidc

Use this data source to retrieve a OIDC IdP from Okta.

## Example Usage

```hcl
data "okta_idp_oidc" "example" {
  name = "Example Provider"
}
```

## Arguments Reference

- `name` - (Optional) The name of the idp to retrieve, conflicts with `id`.

- `id` - (Optional) The id of the idp to retrieve, conflicts with `name`.

## Attributes Reference

- `id` - id of idp.

- `name` - name of the idp.

- `type` - type of idp.

- `authorization_url` - IdP Authorization Server (AS) endpoint to request consent from the user and obtain an authorization code grant.
  
- `authorization_binding` - The method of making an authorization request.
  
- `token_url` - IdP Authorization Server (AS) endpoint to exchange the authorization code grant for an access token.
  
- `token_binding` - The method of making a token request.
  
- `user_info_url` - Protected resource endpoint that returns claims about the authenticated user.
  
- `user_info_binding` - The method of making a user info request.
  
- `jwks_url` - Endpoint where the keys signer publishes its keys in a JWK Set.
  
- `jwks_binding` - The method of making a request for the OIDC JWKS.
  
- `scopes` - The scopes of the IdP.
  
- `protocol_type` - The type of protocol to use.
  
- `client_id` - Unique identifier issued by AS for the Okta IdP instance.
  
- `client_secret` - Client secret issued by AS for the Okta IdP instance.
  
- `issuer_url` - URI that identifies the issuer.
  
- `issuer_mode` - Indicates whether Okta uses the original Okta org domain URL, or a custom domain URL.
  
- `max_clock_skew` - Maximum allowable clock-skew when processing messages from the IdP.
