// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type ProvisioningConnectionProfileResource resource

type ProvisioningConnectionProfile struct {
	AuthScheme string `json:"authScheme,omitempty"`
	Token      string `json:"token,omitempty"`
}

func NewProvisioningConnectionProfile() *ProvisioningConnectionProfile {
	return &ProvisioningConnectionProfile{}
}

func (a *ProvisioningConnectionProfile) IsApplicationInstance() bool {
	return true
}
