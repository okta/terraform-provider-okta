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
