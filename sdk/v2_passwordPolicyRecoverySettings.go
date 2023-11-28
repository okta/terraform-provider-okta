// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type PasswordPolicyRecoverySettings struct {
	Factors *PasswordPolicyRecoveryFactors `json:"factors,omitempty"`
}

func NewPasswordPolicyRecoverySettings() *PasswordPolicyRecoverySettings {
	return &PasswordPolicyRecoverySettings{}
}

func (a *PasswordPolicyRecoverySettings) IsPolicyInstance() bool {
	return true
}
