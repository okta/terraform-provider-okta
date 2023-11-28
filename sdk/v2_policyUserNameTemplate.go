// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type PolicyUserNameTemplate struct {
	Template string `json:"template,omitempty"`
}

func NewPolicyUserNameTemplate() *PolicyUserNameTemplate {
	return &PolicyUserNameTemplate{}
}

func (a *PolicyUserNameTemplate) IsPolicyInstance() bool {
	return true
}
