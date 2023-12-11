// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type ProvisioningConnectionResource resource

type ProvisioningConnection struct {
	Links      interface{} `json:"_links,omitempty"`
	AuthScheme string      `json:"authScheme,omitempty"`
	Status     string      `json:"status,omitempty"`
}

func NewProvisioningConnection() *ProvisioningConnection {
	return &ProvisioningConnection{}
}

func (a *ProvisioningConnection) IsApplicationInstance() bool {
	return true
}
