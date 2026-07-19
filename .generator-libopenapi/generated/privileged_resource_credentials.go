// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PrivilegedResourceCredentials represents the PrivilegedResourceCredentials schema
// Credentials for the privileged resource
type PrivilegedResourceCredentials struct {
	// The password associated with the privileged resource
	Password string `json:"password,omitempty"`
	// The username associated with the privileged resource
	UserName string `json:"userName"`
}
