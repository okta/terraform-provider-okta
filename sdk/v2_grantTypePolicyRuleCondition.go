// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
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
