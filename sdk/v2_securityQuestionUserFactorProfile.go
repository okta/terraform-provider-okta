// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type SecurityQuestionUserFactorProfile struct {
	Answer       string `json:"answer,omitempty"`
	Question     string `json:"question,omitempty"`
	QuestionText string `json:"questionText,omitempty"`
}

func NewSecurityQuestionUserFactorProfile() *SecurityQuestionUserFactorProfile {
	return &SecurityQuestionUserFactorProfile{}
}

func (a *SecurityQuestionUserFactorProfile) IsUserFactorInstance() bool {
	return true
}
