// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type OpenIdConnectApplicationSettingsClient struct {
	ApplicationType        string                                        `json:"application_type,omitempty"`
	ClientUri              string                                        `json:"client_uri,omitempty"`
	ConsentMethod          string                                        `json:"consent_method,omitempty"`
	GrantTypes             []*OAuthGrantType                             `json:"grant_types,omitempty"`
	IdpInitiatedLogin      *OpenIdConnectApplicationIdpInitiatedLogin    `json:"idp_initiated_login,omitempty"`
	InitiateLoginUri       string                                        `json:"initiate_login_uri,omitempty"`
	IssuerMode             string                                        `json:"issuer_mode,omitempty"`
	Jwks                   *OpenIdConnectApplicationSettingsClientKeys   `json:"jwks,omitempty"`
	LogoUri                string                                        `json:"logo_uri,omitempty"`
	PolicyUri              string                                        `json:"policy_uri,omitempty"`
	PostLogoutRedirectUris []string                                      `json:"post_logout_redirect_uris,omitempty"`
	RedirectUris           []string                                      `json:"redirect_uris,omitempty"`
	RefreshToken           *OpenIdConnectApplicationSettingsRefreshToken `json:"refresh_token,omitempty"`
	ResponseTypes          []*OAuthResponseType                          `json:"response_types,omitempty"`
	TosUri                 string                                        `json:"tos_uri,omitempty"`
	WildcardRedirect       string                                        `json:"wildcard_redirect,omitempty"`
	JwksUri                string                                        `json:"jwks_uri,omitempty"`
}
