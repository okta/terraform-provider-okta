// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type PolicyAccountLink struct {
	Action string                   `json:"action,omitempty"`
	Filter *PolicyAccountLinkFilter `json:"filter,omitempty"`
}

func NewPolicyAccountLink() *PolicyAccountLink {
	return &PolicyAccountLink{}
}

func (a *PolicyAccountLink) IsPolicyInstance() bool {
	return true
}
