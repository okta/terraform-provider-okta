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
