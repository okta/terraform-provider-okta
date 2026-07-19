// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PrivilegedResourceUpdateRequest represents the PrivilegedResourceUpdateRequest schema
// Update request for a privileged resource
type PrivilegedResourceUpdateRequest struct {
	Profile interface{} `json:"profile,omitempty"`
	// The username associated with the privileged resource
	UserName string `json:"userName,omitempty"`
}
