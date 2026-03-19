// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type ContextPolicyRuleCondition struct {
	Expression string `json:"expression,omitempty"`
}

func NewContextPolicyRuleCondition() *ContextPolicyRuleCondition {
	return &ContextPolicyRuleCondition{}
}

func (a *ContextPolicyRuleCondition) IsPolicyInstance() bool {
	return true
}
