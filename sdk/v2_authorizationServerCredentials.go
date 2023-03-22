package sdk

type AuthorizationServerCredentials struct {
	Signing *AuthorizationServerCredentialsSigningConfig `json:"signing,omitempty"`
}
