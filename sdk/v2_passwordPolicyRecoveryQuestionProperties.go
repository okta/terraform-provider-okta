// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type PasswordPolicyRecoveryQuestionProperties struct {
	Complexity *PasswordPolicyRecoveryQuestionComplexity `json:"complexity,omitempty"`
}

func NewPasswordPolicyRecoveryQuestionProperties() *PasswordPolicyRecoveryQuestionProperties {
	return &PasswordPolicyRecoveryQuestionProperties{}
}

func (a *PasswordPolicyRecoveryQuestionProperties) IsPolicyInstance() bool {
	return true
}
