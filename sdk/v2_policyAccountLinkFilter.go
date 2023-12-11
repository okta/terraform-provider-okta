// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type PolicyAccountLinkFilter struct {
	Groups *PolicyAccountLinkFilterGroups `json:"groups,omitempty"`
}

func NewPolicyAccountLinkFilter() *PolicyAccountLinkFilter {
	return &PolicyAccountLinkFilter{}
}

func (a *PolicyAccountLinkFilter) IsPolicyInstance() bool {
	return true
}
