// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PasswordCredential represents the PasswordCredential schema
// Specifies a password for a user.  When a user has a valid password, imported hashed password, or password hook, and a response object contains a password credential, then the password object is a b...
type PasswordCredential struct {
	// Specifies the password for a user. The password policy validates this password.
	Value string `json:"value,omitempty"`
	Hash interface{} `json:"hash,omitempty"`
	Hook interface{} `json:"hook,omitempty"`
}
