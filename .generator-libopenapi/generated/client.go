// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// Client represents the Client schema
type Client struct {
	// Unique key for the client application. The `client_id` is immutable. When you create a client Application, you can't specify the `client_id` because Okta uses the application ID for the `client_id`.
	ClientId string `json:"client_id,omitempty"`
	// Time at which the `client_secret` expires or 0 if it doesn't expire (measured in unix seconds)
	ClientSecretExpiresAt int `json:"client_secret_expires_at,omitempty"`
	// URL where Okta sends the logout request
	FrontchannelLogoutUri string `json:"frontchannel_logout_uri,omitempty"`
	// Array of redirection URI strings for use in redirect-based flows. All redirect URIs must be absolute URIs and must not include a fragment component. At least one redirect URI and response type is r...
	RedirectUris []string `json:"redirect_uris,omitempty"`
	TokenEndpointAuthMethod interface{} `json:"token_endpoint_auth_method,omitempty"`
	ApplicationType interface{} `json:"application_type,omitempty"`
	// Time at which the `client_id` was issued (measured in unix seconds)
	ClientIdIssuedAt int `json:"client_id_issued_at,omitempty"`
	// Human-readable string name of the client application
	ClientName string `json:"client_name,omitempty"`
	// Include user session details
	FrontchannelLogoutSessionRequired bool `json:"frontchannel_logout_session_required,omitempty"`
	// URL string that references a [JSON Web Key Set](https://tools.ietf.org/html/rfc7517#section-5) for validating JWTs presented to Okta
	JwksUri string `json:"jwks_uri,omitempty"`
	// URL string that references a logo for the client consent dialog (not the sign-in dialog)
	LogoUri string `json:"logo_uri,omitempty"`
	// Array of redirection URI strings for use for relying party initiated logouts
	PostLogoutRedirectUris string `json:"post_logout_redirect_uris,omitempty"`
	// The type of [JSON Web Key Set](https://tools.ietf.org/html/rfc7517#section-5) algorithm that must be used for signing request objects
	RequestObjectSigningAlg []interface{} `json:"request_object_signing_alg,omitempty"`
	// URL string of a web page providing the client's policy document
	PolicyUri string `json:"policy_uri,omitempty"`
	// URL string of a web page providing the client's terms of service document
	TosUri string `json:"tos_uri,omitempty"`
	// OAuth 2.0 client secret string (used for confidential clients). The `client_secret` is shown only on the response of the creation or update of a client Application (and only if the `token_endpoint_...
	ClientSecret string `json:"client_secret,omitempty"`
	// Array of OAuth 2.0 grant type strings. Default value: `[authorization_code]`
	GrantTypes []interface{} `json:"grant_types,omitempty"`
	// URL that a third party can use to initiate a login by the client
	InitiateLoginUri string `json:"initiate_login_uri,omitempty"`
	// Array of OAuth 2.0 response type strings. Default value: `[code]`
	ResponseTypes []interface{} `json:"response_types,omitempty"`
}
