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
