// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
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
