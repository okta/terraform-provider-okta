// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type PolicyRuleActionsEnroll struct {
	Self string `json:"self,omitempty"`
}

func NewPolicyRuleActionsEnroll() *PolicyRuleActionsEnroll {
	return &PolicyRuleActionsEnroll{}
}

func (a *PolicyRuleActionsEnroll) IsPolicyInstance() bool {
	return true
}
