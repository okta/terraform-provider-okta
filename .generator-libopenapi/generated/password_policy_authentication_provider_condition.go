// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PasswordPolicyAuthenticationProviderCondition represents the PasswordPolicyAuthenticationProviderCondition schema
// Specifies an authentication provider that's the source of some or all users
type PasswordPolicyAuthenticationProviderCondition struct {
	Provider interface{} `json:"provider,omitempty"`
	Include []string `json:"include,omitempty"`
}
