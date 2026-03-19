// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type RiskPolicyRuleCondition struct {
	Behaviors []string `json:"behaviors,omitempty"`
}

func NewRiskPolicyRuleCondition() *RiskPolicyRuleCondition {
	return &RiskPolicyRuleCondition{}
}

func (a *RiskPolicyRuleCondition) IsPolicyInstance() bool {
	return true
}
