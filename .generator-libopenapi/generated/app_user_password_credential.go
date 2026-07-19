// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AppUserPasswordCredential represents the AppUserPasswordCredential schema
// The user's password. This is a write-only property. An empty `password` object is returned to indicate that a password value exists.
type AppUserPasswordCredential struct {
	// Password value
	Value string `json:"value,omitempty"`
}
