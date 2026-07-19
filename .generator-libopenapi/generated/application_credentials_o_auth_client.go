// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ApplicationCredentialsOAuthClient represents the ApplicationCredentialsOAuthClient schema
type ApplicationCredentialsOAuthClient struct {
	TokenEndpointAuthMethod interface{} `json:"token_endpoint_auth_method,omitempty"`
	// Requested key rotation mode
	AutoKeyRotation bool `json:"autoKeyRotation,omitempty"`
	// Unique identifier for the OAuth 2.0 client app  > **Notes:** > * If you don't specify the `client_id`, this immutable property is populated with the [Application instance ID](/openapi/okta-manageme...
	ClientId string `json:"client_id,omitempty"`
	// OAuth 2.0 client secret string (used for confidential clients)  > **Notes:** If a `client_secret` isn't provided on creation, and the `token_endpoint_auth_method` requires one, Okta generates a ran...
	ClientSecret string `json:"client_secret,omitempty"`
	// Requires Proof Key for Code Exchange (PKCE) for additional verification. If `token_endpoint_auth_method` is `none`, then `pkce_required` must be `true`. The default is `true` for browser and native...
	PkceRequired bool `json:"pkce_required,omitempty"`
}
