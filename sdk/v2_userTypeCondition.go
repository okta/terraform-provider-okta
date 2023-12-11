// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type UserTypeCondition struct {
	Exclude []string `json:"exclude,omitempty"`
	Include []string `json:"include,omitempty"`
}

func NewUserTypeCondition() *UserTypeCondition {
	return &UserTypeCondition{}
}

func (a *UserTypeCondition) IsPolicyInstance() bool {
	return true
}
