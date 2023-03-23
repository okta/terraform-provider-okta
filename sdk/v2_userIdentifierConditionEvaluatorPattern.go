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
