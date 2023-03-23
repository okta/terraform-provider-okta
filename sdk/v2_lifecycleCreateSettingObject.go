package sdk

type LifecycleCreateSettingObject struct {
	Status string `json:"status,omitempty"`
}

func NewLifecycleCreateSettingObject() *LifecycleCreateSettingObject {
	return &LifecycleCreateSettingObject{}
}

func (a *LifecycleCreateSettingObject) IsApplicationInstance() bool {
	return true
}
