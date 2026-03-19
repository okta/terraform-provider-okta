// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type UserCredentials struct {
	Password         *PasswordCredential         `json:"password,omitempty"`
	Provider         *AuthenticationProvider     `json:"provider,omitempty"`
	RecoveryQuestion *RecoveryQuestionCredential `json:"recovery_question,omitempty"`
}
