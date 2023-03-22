package sdk

type AccessPolicyConstraints struct {
	Knowledge  *KnowledgeConstraint  `json:"knowledge,omitempty"`
	Possession *PossessionConstraint `json:"possession,omitempty"`
}

func NewAccessPolicyConstraints() *AccessPolicyConstraints {
	return &AccessPolicyConstraints{}
}

func (a *AccessPolicyConstraints) IsPolicyInstance() bool {
	return true
}
