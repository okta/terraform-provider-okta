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
