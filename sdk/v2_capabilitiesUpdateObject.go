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
