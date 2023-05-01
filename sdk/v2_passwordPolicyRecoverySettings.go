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
