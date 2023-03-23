package sdk

type PolicyPeopleCondition struct {
	Groups *GroupCondition `json:"groups,omitempty"`
	Users  *UserCondition  `json:"users,omitempty"`
}

func NewPolicyPeopleCondition() *PolicyPeopleCondition {
	return &PolicyPeopleCondition{}
}

func (a *PolicyPeopleCondition) IsPolicyInstance() bool {
	return true
}
