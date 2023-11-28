// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
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
