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
