// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
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
