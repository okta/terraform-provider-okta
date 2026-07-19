// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserCredentialsWritable represents the UserCredentialsWritable schema
// Specifies primary authentication and recovery credentials for a user. Credential types and requirements vary depending on the provider and security policy of the org.
type UserCredentialsWritable struct {
	Password interface{} `json:"password,omitempty"`
	Provider interface{} `json:"provider,omitempty"`
	RecoveryQuestion interface{} `json:"recovery_question,omitempty"`
}
