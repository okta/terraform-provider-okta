// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type PasswordPolicyRecoveryEmailProperties struct {
	RecoveryToken *PasswordPolicyRecoveryEmailRecoveryToken `json:"recoveryToken,omitempty"`
}

func NewPasswordPolicyRecoveryEmailProperties() *PasswordPolicyRecoveryEmailProperties {
	return &PasswordPolicyRecoveryEmailProperties{}
}

func (a *PasswordPolicyRecoveryEmailProperties) IsPolicyInstance() bool {
	return true
}
