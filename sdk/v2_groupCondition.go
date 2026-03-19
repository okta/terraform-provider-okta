// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
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
