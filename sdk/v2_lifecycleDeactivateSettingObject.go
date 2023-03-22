package sdk

type LifecycleDeactivateSettingObject struct {
	Status string `json:"status,omitempty"`
}

func NewLifecycleDeactivateSettingObject() *LifecycleDeactivateSettingObject {
	return &LifecycleDeactivateSettingObject{}
}

func (a *LifecycleDeactivateSettingObject) IsApplicationInstance() bool {
	return true
}
