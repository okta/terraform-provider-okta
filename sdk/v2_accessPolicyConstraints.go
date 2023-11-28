// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
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
