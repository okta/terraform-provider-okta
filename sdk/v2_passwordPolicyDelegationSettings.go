// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
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
