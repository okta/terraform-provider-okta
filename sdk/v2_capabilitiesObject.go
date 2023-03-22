package sdk

type CapabilitiesObject struct {
	Create *CapabilitiesCreateObject `json:"create,omitempty"`
	Update *CapabilitiesUpdateObject `json:"update,omitempty"`
}

func NewCapabilitiesObject() *CapabilitiesObject {
	return &CapabilitiesObject{}
}

func (a *CapabilitiesObject) IsApplicationInstance() bool {
	return true
}
