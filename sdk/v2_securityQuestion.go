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
