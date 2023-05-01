package sdk

type AuthenticatorProvider struct {
	Configuration *AuthenticatorProviderConfiguration `json:"configuration,omitempty"`
	Type          string                              `json:"type,omitempty"`
}
