package sdk

type GroupCondition struct {
	Exclude []string `json:"exclude,omitempty"`
	Include []string `json:"include,omitempty"`
}

func NewGroupCondition() *GroupCondition {
	return &GroupCondition{}
}

func (a *GroupCondition) IsPolicyInstance() bool {
	return true
}
