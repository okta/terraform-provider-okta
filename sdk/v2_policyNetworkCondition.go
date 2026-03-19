// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type PolicyNetworkCondition struct {
	Connection string   `json:"connection,omitempty"`
	Exclude    []string `json:"exclude,omitempty"`
	Include    []string `json:"include,omitempty"`
}

func NewPolicyNetworkCondition() *PolicyNetworkCondition {
	return &PolicyNetworkCondition{}
}

func (a *PolicyNetworkCondition) IsPolicyInstance() bool {
	return true
}
