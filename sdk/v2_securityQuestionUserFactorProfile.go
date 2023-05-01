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
