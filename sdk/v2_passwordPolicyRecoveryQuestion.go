package sdk

type PasswordPolicyRecoveryQuestion struct {
	Properties *PasswordPolicyRecoveryQuestionProperties `json:"properties,omitempty"`
	Status     string                                    `json:"status,omitempty"`
}

func NewPasswordPolicyRecoveryQuestion() *PasswordPolicyRecoveryQuestion {
	return &PasswordPolicyRecoveryQuestion{}
}

func (a *PasswordPolicyRecoveryQuestion) IsPolicyInstance() bool {
	return true
}
