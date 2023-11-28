// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
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
