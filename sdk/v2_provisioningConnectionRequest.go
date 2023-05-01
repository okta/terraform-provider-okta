package sdk

type ProvisioningConnectionRequest struct {
	Profile *ProvisioningConnectionProfile `json:"profile,omitempty"`
}

func NewProvisioningConnectionRequest() *ProvisioningConnectionRequest {
	return &ProvisioningConnectionRequest{}
}

func (a *ProvisioningConnectionRequest) IsApplicationInstance() bool {
	return true
}
