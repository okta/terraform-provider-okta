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
