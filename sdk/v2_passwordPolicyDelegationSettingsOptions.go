// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
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
