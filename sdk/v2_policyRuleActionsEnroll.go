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
