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
