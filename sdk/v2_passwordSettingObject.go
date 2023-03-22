package sdk

type PasswordSettingObject struct {
	Change string `json:"change,omitempty"`
	Seed   string `json:"seed,omitempty"`
	Status string `json:"status,omitempty"`
}

func NewPasswordSettingObject() *PasswordSettingObject {
	return &PasswordSettingObject{}
}

func (a *PasswordSettingObject) IsApplicationInstance() bool {
	return true
}
