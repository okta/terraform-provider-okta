package sdk

type ClientSecretMetadata struct {
	ClientSecret string `json:"client_secret,omitempty"`
}

func NewClientSecretMetadata() *ClientSecretMetadata {
	return &ClientSecretMetadata{}
}

func (a *ClientSecretMetadata) IsApplicationInstance() bool {
	return true
}
