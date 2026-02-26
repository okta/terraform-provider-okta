// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

// Office365ClientCondition represents the Office 365 client condition for access policy rules.
// This condition allows filtering based on Office 365 client types such as WEB, MODERN_AUTH, AAD_JOIN, etc.
type Office365ClientCondition struct {
	Include []string `json:"include,omitempty"`
}

func NewOffice365ClientCondition() *Office365ClientCondition {
	return &Office365ClientCondition{}
}

func (a *Office365ClientCondition) IsPolicyInstance() bool {
	return true
}
