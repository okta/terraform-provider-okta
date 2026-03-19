// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
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
