// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
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
