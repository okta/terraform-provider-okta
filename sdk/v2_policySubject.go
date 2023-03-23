package sdk

type PolicySubject struct {
	Filter           string                  `json:"filter,omitempty"`
	Format           []string                `json:"format,omitempty"`
	MatchAttribute   string                  `json:"matchAttribute,omitempty"`
	MatchType        string                  `json:"matchType,omitempty"`
	UserNameTemplate *PolicyUserNameTemplate `json:"userNameTemplate,omitempty"`
}

func NewPolicySubject() *PolicySubject {
	return &PolicySubject{}
}

func (a *PolicySubject) IsPolicyInstance() bool {
	return true
}
