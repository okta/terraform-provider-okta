// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type SecurityQuestion struct {
	Answer       string `json:"answer,omitempty"`
	Question     string `json:"question,omitempty"`
	QuestionText string `json:"questionText,omitempty"`
}

func NewSecurityQuestion() *SecurityQuestion {
	return &SecurityQuestion{}
}

func (a *SecurityQuestion) IsUserFactorInstance() bool {
	return true
}
