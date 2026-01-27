// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type KnowledgeConstraint struct {
	AuthenticationMethods         []AuthenticationMethodObject `json:"authenticationMethods,omitempty"`
	ExcludedAuthenticationMethods []AuthenticationMethodObject `json:"excludedAuthenticationMethods,omitempty"`
	Methods                       []string                     `json:"methods,omitempty"`
	ReauthenticateIn              string                       `json:"reauthenticateIn,omitempty"`
	Types                         []string                     `json:"types,omitempty"`
	Required                      bool                         `json:"required"`
}

type AuthenticationMethodObject struct {
	Key    string `json:"key,omitempty"`
	Id     string `json:"id,omitempty"`
	Method string `json:"method,omitempty"`
}

func NewKnowledgeConstraint() *KnowledgeConstraint {
	return &KnowledgeConstraint{}
}

func (a *KnowledgeConstraint) IsPolicyInstance() bool {
	return true
}
