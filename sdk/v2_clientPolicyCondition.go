package sdk

type ClientPolicyCondition struct {
	Include []string `json:"include,omitempty"`
}

func NewClientPolicyCondition() *ClientPolicyCondition {
	return &ClientPolicyCondition{}
}

func (a *ClientPolicyCondition) IsPolicyInstance() bool {
	return true
}
