package sdk

type PasswordPolicyAuthenticationProviderCondition struct {
	Include  []string `json:"include,omitempty"`
	Provider string   `json:"provider,omitempty"`
}

func NewPasswordPolicyAuthenticationProviderCondition() *PasswordPolicyAuthenticationProviderCondition {
	return &PasswordPolicyAuthenticationProviderCondition{}
}

func (a *PasswordPolicyAuthenticationProviderCondition) IsPolicyInstance() bool {
	return true
}
