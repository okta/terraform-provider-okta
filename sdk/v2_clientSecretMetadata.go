// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
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
