// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type PolicyRuleAuthContextCondition struct {
	AuthType string `json:"authType,omitempty"`
}

func NewPolicyRuleAuthContextCondition() *PolicyRuleAuthContextCondition {
	return &PolicyRuleAuthContextCondition{}
}

func (a *PolicyRuleAuthContextCondition) IsPolicyInstance() bool {
	return true
}
