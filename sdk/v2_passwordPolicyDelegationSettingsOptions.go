package sdk

type PasswordPolicyDelegationSettingsOptions struct {
	SkipUnlock *bool `json:"skipUnlock,omitempty"`
}

func NewPasswordPolicyDelegationSettingsOptions() *PasswordPolicyDelegationSettingsOptions {
	return &PasswordPolicyDelegationSettingsOptions{}
}

func (a *PasswordPolicyDelegationSettingsOptions) IsPolicyInstance() bool {
	return true
}
