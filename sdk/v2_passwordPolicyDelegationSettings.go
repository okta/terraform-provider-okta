package sdk

type PasswordPolicyDelegationSettings struct {
	Options *PasswordPolicyDelegationSettingsOptions `json:"options,omitempty"`
}

func NewPasswordPolicyDelegationSettings() *PasswordPolicyDelegationSettings {
	return &PasswordPolicyDelegationSettings{}
}

func (a *PasswordPolicyDelegationSettings) IsPolicyInstance() bool {
	return true
}
