package sdk

type UserCredentials struct {
	Password         *PasswordCredential         `json:"password,omitempty"`
	Provider         *AuthenticationProvider     `json:"provider,omitempty"`
	RecoveryQuestion *RecoveryQuestionCredential `json:"recovery_question,omitempty"`
}
