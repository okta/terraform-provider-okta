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
