package sdk

type IdentityProviderCredentials struct {
	Client  *IdentityProviderCredentialsClient  `json:"client,omitempty"`
	Signing *IdentityProviderCredentialsSigning `json:"signing,omitempty"`
	Trust   *IdentityProviderCredentialsTrust   `json:"trust,omitempty"`
}
