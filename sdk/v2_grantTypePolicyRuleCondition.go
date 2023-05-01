package sdk

type GrantTypePolicyRuleCondition struct {
	Include []string `json:"include,omitempty"`
}

func NewGrantTypePolicyRuleCondition() *GrantTypePolicyRuleCondition {
	return &GrantTypePolicyRuleCondition{}
}

func (a *GrantTypePolicyRuleCondition) IsPolicyInstance() bool {
	return true
}
