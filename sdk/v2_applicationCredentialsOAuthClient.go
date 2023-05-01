package sdk

type ApplicationCredentialsOAuthClient struct {
	AutoKeyRotation         *bool  `json:"autoKeyRotation,omitempty"`
	ClientId                string `json:"client_id,omitempty"`
	ClientSecret            string `json:"client_secret,omitempty"`
	PkceRequired            *bool  `json:"pkce_required,omitempty"`
	TokenEndpointAuthMethod string `json:"token_endpoint_auth_method,omitempty"`
}
