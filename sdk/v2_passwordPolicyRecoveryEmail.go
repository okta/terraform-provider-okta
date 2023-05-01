package sdk

type PasswordPolicyRecoveryEmail struct {
	Properties *PasswordPolicyRecoveryEmailProperties `json:"properties,omitempty"`
	Status     string                                 `json:"status,omitempty"`
}

func NewPasswordPolicyRecoveryEmail() *PasswordPolicyRecoveryEmail {
	return &PasswordPolicyRecoveryEmail{}
}

func (a *PasswordPolicyRecoveryEmail) IsPolicyInstance() bool {
	return true
}
