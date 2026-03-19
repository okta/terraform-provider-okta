// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type PasswordPolicyRecoveryFactors struct {
	OktaCall         *PasswordPolicyRecoveryFactorSettings `json:"okta_call,omitempty"`
	OktaEmail        *PasswordPolicyRecoveryEmail          `json:"okta_email,omitempty"`
	OktaSms          *PasswordPolicyRecoveryFactorSettings `json:"okta_sms,omitempty"`
	RecoveryQuestion *PasswordPolicyRecoveryQuestion       `json:"recovery_question,omitempty"`
}

func NewPasswordPolicyRecoveryFactors() *PasswordPolicyRecoveryFactors {
	return &PasswordPolicyRecoveryFactors{}
}

func (a *PasswordPolicyRecoveryFactors) IsPolicyInstance() bool {
	return true
}
