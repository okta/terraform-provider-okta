// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type AccessPolicyRuleCustomCondition struct {
	Condition string `json:"condition,omitempty"`
}

func NewAccessPolicyRuleCustomCondition() *AccessPolicyRuleCustomCondition {
	return &AccessPolicyRuleCustomCondition{}
}

func (a *AccessPolicyRuleCustomCondition) IsPolicyInstance() bool {
	return true
}
