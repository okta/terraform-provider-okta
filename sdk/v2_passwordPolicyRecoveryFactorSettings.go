// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type PasswordPolicyRecoveryFactorSettings struct {
	Status string `json:"status,omitempty"`
}

func NewPasswordPolicyRecoveryFactorSettings() *PasswordPolicyRecoveryFactorSettings {
	return &PasswordPolicyRecoveryFactorSettings{
		Status: "INACTIVE",
	}
}

func (a *PasswordPolicyRecoveryFactorSettings) IsPolicyInstance() bool {
	return true
}
