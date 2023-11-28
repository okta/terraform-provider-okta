// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type AppAndInstanceConditionEvaluatorAppOrInstance struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

func NewAppAndInstanceConditionEvaluatorAppOrInstance() *AppAndInstanceConditionEvaluatorAppOrInstance {
	return &AppAndInstanceConditionEvaluatorAppOrInstance{}
}

func (a *AppAndInstanceConditionEvaluatorAppOrInstance) IsPolicyInstance() bool {
	return true
}
