package sdk

type PasswordPolicyPasswordSettings struct {
	Age        *PasswordPolicyPasswordSettingsAge        `json:"age,omitempty"`
	Complexity *PasswordPolicyPasswordSettingsComplexity `json:"complexity,omitempty"`
	Lockout    *PasswordPolicyPasswordSettingsLockout    `json:"lockout,omitempty"`
}

func NewPasswordPolicyPasswordSettings() *PasswordPolicyPasswordSettings {
	return &PasswordPolicyPasswordSettings{}
}

func (a *PasswordPolicyPasswordSettings) IsPolicyInstance() bool {
	return true
}
