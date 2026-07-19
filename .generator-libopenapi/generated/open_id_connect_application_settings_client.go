// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OpenIdConnectApplicationSettingsClient represents the OpenIdConnectApplicationSettingsClient schema
type OpenIdConnectApplicationSettingsClient struct {
	// The URL string that references a logo for the client. This logo appears on the client tile in the End-User Dashboard. It also appears on the client consent dialog during the client consent flow.
	LogoUri string `json:"logo_uri,omitempty"`
	// URL string of a web page providing the client's terms of service document
	TosUri string `json:"tos_uri,omitempty"`
	RefreshToken interface{} `json:"refresh_token,omitempty"`
	// Type of the subject
	SubjectType string `json:"subject_type,omitempty"`
	// The ID of the custom authenticator that authenticates the user > **Note:** This property appears for clients with `urn:openid:params:grant-type:ciba` defined as one of the `grant_types`.
	BackchannelCustomAuthenticatorId string `json:"backchannel_custom_authenticator_id,omitempty"`
	// The delivery mode for Client-Initiated Backchannel Authentication (CIBA).  Currently, only `poll` is supported. > **Note:** This property appears for clients with `urn:openid:params:grant-type:ciba...
	BackchannelTokenDeliveryMode string `json:"backchannel_token_delivery_mode,omitempty"`
	GrantTypes []interface{} `json:"grant_types"`
	// URL string that a third party can use to initiate the sign-in flow by the client
	InitiateLoginUri string `json:"initiate_login_uri,omitempty"`
	// URL string that references a JSON Web Key Set for validating JWTs presented to Okta or for encrypting ID tokens minted by Okta for the client
	JwksUri string `json:"jwks_uri,omitempty"`
	// URL string of a web page providing the client's policy document
	PolicyUri string `json:"policy_uri,omitempty"`
	// Array of redirection URI strings for relying party-initiated logouts
	PostLogoutRedirectUris []string `json:"post_logout_redirect_uris,omitempty"`
	IdTokenEncryptedResponseAlg interface{} `json:"id_token_encrypted_response_alg,omitempty"`
	Jwks interface{} `json:"jwks,omitempty"`
	// The signing algorithm for Client-Initiated Backchannel Authentication (CIBA) signed requests using JWT. If this value isn't set and a JWT-signed request is sent, the request fails. > **Note:** This...
	BackchannelAuthenticationRequestSigningAlg string `json:"backchannel_authentication_request_signing_alg,omitempty"`
	IdpInitiatedLogin interface{} `json:"idp_initiated_login,omitempty"`
	ApplicationType interface{} `json:"application_type,omitempty"`
	ConsentMethod interface{} `json:"consent_method,omitempty"`
	// <x-lifecycle-container><x-lifecycle class="ea"></x-lifecycle> <x-lifecycle class="oie"></x-lifecycle></x-lifecycle-container>Determines whether Okta sends `sid` and `iss` in the logout request
	FrontchannelLogoutSessionRequired bool `json:"frontchannel_logout_session_required,omitempty"`
	Network interface{} `json:"network,omitempty"`
	// <x-lifecycle-container><x-lifecycle class="ea"></x-lifecycle> <x-lifecycle class="oie"></x-lifecycle></x-lifecycle-container>Allows the app to participate in front-channel Single Logout  > **Note:*...
	ParticipateSlo bool `json:"participate_slo,omitempty"`
	// Array of redirection URI strings for use in redirect-based flows. > **Note:** At least one `redirect_uris` and `response_types` are required for all client types, with exceptions: if the client use...
	RedirectUris []string `json:"redirect_uris,omitempty"`
	// The sector identifier used for pairwise `subject_type`. See [OIDC Pairwise Identifier Algorithm](https://openid.net/specs/openid-connect-messages-1_0-20.html#idtype.pairwise.alg)
	SectorIdentifierUri string `json:"sector_identifier_uri,omitempty"`
	// Indicates if the client is allowed to use wildcard matching of `redirect_uris`
	WildcardRedirect string `json:"wildcard_redirect,omitempty"`
	// URL string of a web page providing information about the client
	ClientUri string `json:"client_uri,omitempty"`
	// Indicates that the client application uses Demonstrating Proof-of-Possession (DPoP) for token requests. If `true`, the authorization server rejects token requests from this client that don't contai...
	DpopBoundAccessTokens bool `json:"dpop_bound_access_tokens,omitempty"`
	IssuerMode interface{} `json:"issuer_mode,omitempty"`
	// The type of JSON Web Key Set (JWKS) algorithm that must be used for signing request objects
	RequestObjectSigningAlg string `json:"request_object_signing_alg,omitempty"`
	// Array of OAuth 2.0 response type strings
	ResponseTypes []interface{} `json:"response_types,omitempty"`
	// <x-lifecycle-container><x-lifecycle class="ea"></x-lifecycle> <x-lifecycle class="oie"></x-lifecycle></x-lifecycle-container>URL where Okta sends the logout request
	FrontchannelLogoutUri string `json:"frontchannel_logout_uri,omitempty"`
}
