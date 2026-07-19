// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OAuthCredentialsClient represents the OAuthCredentialsClient schema
// OAuth 2.0 and OpenID Connect Client object > **Note:** You must complete client registration with the IdP Authorization Server for your Okta IdP instance to obtain client credentials.
type OAuthCredentialsClient struct {
	// The [Unique identifier](https://tools.ietf.org/html/rfc6749#section-2.2) issued by the AS for the Okta IdP instance
	ClientId string `json:"client_id,omitempty"`
	// The [Client secret](https://tools.ietf.org/html/rfc6749#section-2.3.1) issued by the AS for the Okta IdP instance
	ClientSecret string `json:"client_secret,omitempty"`
	// Require Proof Key for Code Exchange (PKCE) for additional verification
	PkceRequired bool `json:"pkce_required,omitempty"`
	// Client authentication methods supported by the token endpoint
	TokenEndpointAuthMethod string `json:"token_endpoint_auth_method,omitempty"`
}
