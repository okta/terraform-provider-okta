// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type CapabilitiesCreateObject struct {
	LifecycleCreate *LifecycleCreateSettingObject `json:"lifecycleCreate,omitempty"`
}

func NewCapabilitiesCreateObject() *CapabilitiesCreateObject {
	return &CapabilitiesCreateObject{}
}

func (a *CapabilitiesCreateObject) IsApplicationInstance() bool {
	return true
}
