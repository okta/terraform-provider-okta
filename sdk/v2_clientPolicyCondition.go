// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
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
