// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OAuthMetadata represents the OAuthMetadata schema
type OAuthMetadata struct {
	// URL of the authorization server's authorization endpoint.
	AuthorizationEndpoint string `json:"authorization_endpoint,omitempty"`
	// <x-lifecycle-container><x-lifecycle class="oie"></x-lifecycle></x-lifecycle-container>A list of signing algorithms that this authorization server supports for signed requests.
	BackchannelAuthenticationRequestSigningAlgValuesSupported []interface{} `json:"backchannel_authentication_request_signing_alg_values_supported,omitempty"`
	// A list of PKCE code challenge methods supported by this authorization server.
	CodeChallengeMethodsSupported []interface{} `json:"code_challenge_methods_supported,omitempty"`
	// A list of client authentication methods supported by this introspection endpoint.
	IntrospectionEndpointAuthMethodsSupported []interface{} `json:"introspection_endpoint_auth_methods_supported,omitempty"`
	// The authorization server's issuer identifier. In the context of this document, this is your authorization server's base URL. This becomes the `iss` claim in an access token.
	Issuer string `json:"issuer,omitempty"`
	// URL of the authorization server's JSON Web Key Set document.
	JwksUri string `json:"jwks_uri,omitempty"`
	// A list of the `response_mode` values that this authorization server supports. More information here.
	ResponseModesSupported []interface{} `json:"response_modes_supported,omitempty"`
	// URL of the authorization server's revocation endpoint.
	RevocationEndpoint string `json:"revocation_endpoint,omitempty"`
	// A list of the claims supported by this authorization server.
	ClaimsSupported []interface{} `json:"claims_supported,omitempty"`
	DeviceAuthorizationEndpoint string `json:"device_authorization_endpoint,omitempty"`
	// A list of signing algorithms supported by this authorization server for Demonstrating Proof-of-Possession (DPoP) JWTs.
	DpopSigningAlgValuesSupported []string `json:"dpop_signing_alg_values_supported,omitempty"`
	// A list of signing algorithms that this authorization server supports for signed requests.
	RequestObjectSigningAlgValuesSupported []interface{} `json:"request_object_signing_alg_values_supported,omitempty"`
	// Indicates if Request Parameters are supported by this authorization server.
	RequestParameterSupported bool `json:"request_parameter_supported,omitempty"`
	// A list of the `response_type` values that this authorization server supports. Can be a combination of `code`, `token`, and `id_token`.
	ResponseTypesSupported []interface{} `json:"response_types_supported,omitempty"`
	// A list of client authentication methods supported by this revocation endpoint.
	RevocationEndpointAuthMethodsSupported []interface{} `json:"revocation_endpoint_auth_methods_supported,omitempty"`
	// A list of the Subject Identifier types that this authorization server supports. Valid types include `pairwise` and `public`, but only `public` is currently supported. See the [Subject Identifier Ty...
	SubjectTypesSupported []interface{} `json:"subject_types_supported,omitempty"`
	PushedAuthorizationRequestEndpoint string `json:"pushed_authorization_request_endpoint,omitempty"`
	// <x-lifecycle-container><x-lifecycle class="oie"></x-lifecycle></x-lifecycle-container>The delivery modes that this authorization server supports for Client-Initiated Backchannel Authentication.
	BackchannelTokenDeliveryModesSupported []interface{} `json:"backchannel_token_delivery_modes_supported,omitempty"`
	// A list of the grant type values that this authorization server supports.
	GrantTypesSupported []interface{} `json:"grant_types_supported,omitempty"`
	// URL of the authorization server's JSON Web Key Set document.
	RegistrationEndpoint string `json:"registration_endpoint,omitempty"`
	// A list of the scope values that this authorization server supports.
	ScopesSupported []interface{} `json:"scopes_supported,omitempty"`
	// A list of client authentication methods supported by this token endpoint.
	TokenEndpointAuthMethodsSupported []interface{} `json:"token_endpoint_auth_methods_supported,omitempty"`
	// URL of the authorization server's logout endpoint.
	EndSessionEndpoint string `json:"end_session_endpoint,omitempty"`
	// URL of the authorization server's introspection endpoint.
	IntrospectionEndpoint string `json:"introspection_endpoint,omitempty"`
	// URL of the authorization server's token endpoint.
	TokenEndpoint string `json:"token_endpoint,omitempty"`
}
