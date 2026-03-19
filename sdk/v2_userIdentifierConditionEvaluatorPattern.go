// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type UserIdentifierConditionEvaluatorPattern struct {
	MatchType string `json:"matchType,omitempty"`
	Value     string `json:"value,omitempty"`
}

func NewUserIdentifierConditionEvaluatorPattern() *UserIdentifierConditionEvaluatorPattern {
	return &UserIdentifierConditionEvaluatorPattern{}
}

func (a *UserIdentifierConditionEvaluatorPattern) IsPolicyInstance() bool {
	return true
}
