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
