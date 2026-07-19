// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserCredentials represents the UserCredentials schema
// Specifies primary authentication and recovery credentials for a user. Credential types and requirements vary depending on the provider and security policy of the org.
type UserCredentials struct {
	Provider interface{} `json:"provider,omitempty"`
	RecoveryQuestion interface{} `json:"recovery_question,omitempty"`
	Password interface{} `json:"password,omitempty"`
}
