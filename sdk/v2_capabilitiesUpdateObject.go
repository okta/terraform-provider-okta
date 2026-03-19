// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type CapabilitiesUpdateObject struct {
	LifecycleDeactivate *LifecycleDeactivateSettingObject `json:"lifecycleDeactivate,omitempty"`
	Password            *PasswordSettingObject            `json:"password,omitempty"`
	Profile             *ProfileSettingObject             `json:"profile,omitempty"`
}

func NewCapabilitiesUpdateObject() *CapabilitiesUpdateObject {
	return &CapabilitiesUpdateObject{}
}

func (a *CapabilitiesUpdateObject) IsApplicationInstance() bool {
	return true
}
