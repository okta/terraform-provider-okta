// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type PolicyAccountLinkFilterGroups struct {
	Include []string `json:"include,omitempty"`
}

func NewPolicyAccountLinkFilterGroups() *PolicyAccountLinkFilterGroups {
	return &PolicyAccountLinkFilterGroups{}
}

func (a *PolicyAccountLinkFilterGroups) IsPolicyInstance() bool {
	return true
}
