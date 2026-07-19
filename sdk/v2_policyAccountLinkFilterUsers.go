// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type PolicyAccountLinkFilterUsers struct {
	Exclude       []string `json:"exclude,omitempty"`
	ExcludeAdmins *bool    `json:"excludeAdmins,omitempty"`
}

func NewPolicyAccountLinkFilterUsers() *PolicyAccountLinkFilterUsers {
	return &PolicyAccountLinkFilterUsers{}
}

func (a *PolicyAccountLinkFilterUsers) IsPolicyInstance() bool {
	return true
}