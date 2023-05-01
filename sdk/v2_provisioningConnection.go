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
