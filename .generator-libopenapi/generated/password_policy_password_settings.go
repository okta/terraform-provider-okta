// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PasswordPolicyPasswordSettings represents the PasswordPolicyPasswordSettings schema
// Specifies the password settings for the policy
type PasswordPolicyPasswordSettings struct {
	Complexity interface{} `json:"complexity,omitempty"`
	Lockout interface{} `json:"lockout,omitempty"`
	BreachedProtection interface{} `json:"breachedProtection,omitempty"`
	Age interface{} `json:"age,omitempty"`
}
