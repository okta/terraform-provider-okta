// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type UserCondition struct {
	Exclude []string `json:"exclude,omitempty"`
	Include []string `json:"include,omitempty"`
}

func NewUserCondition() *UserCondition {
	return &UserCondition{}
}

func (a *UserCondition) IsPolicyInstance() bool {
	return true
}
