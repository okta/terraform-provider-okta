// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type ProfileSettingObject struct {
	Status string `json:"status,omitempty"`
}

func NewProfileSettingObject() *ProfileSettingObject {
	return &ProfileSettingObject{}
}

func (a *ProfileSettingObject) IsApplicationInstance() bool {
	return true
}
