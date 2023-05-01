package sdk

type ApplicationFeatureResource resource

type ApplicationFeature struct {
	Links        interface{}         `json:"_links,omitempty"`
	Capabilities *CapabilitiesObject `json:"capabilities,omitempty"`
	Description  string              `json:"description,omitempty"`
	Name         string              `json:"name,omitempty"`
	Status       string              `json:"status,omitempty"`
}

func NewApplicationFeature() *ApplicationFeature {
	return &ApplicationFeature{}
}

func (a *ApplicationFeature) IsApplicationInstance() bool {
	return true
}
