package sdk

type KnowledgeConstraint struct {
	Methods          []string `json:"methods,omitempty"`
	ReauthenticateIn string   `json:"reauthenticateIn,omitempty"`
	Types            []string `json:"types,omitempty"`
}

func NewKnowledgeConstraint() *KnowledgeConstraint {
	return &KnowledgeConstraint{}
}

func (a *KnowledgeConstraint) IsPolicyInstance() bool {
	return true
}
