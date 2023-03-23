package sdk

type IdentityProviderCredentialsSigning struct {
	Kid        string `json:"kid,omitempty"`
	PrivateKey string `json:"privateKey,omitempty"`
	TeamId     string `json:"teamId,omitempty"`
}
