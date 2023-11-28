// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
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
