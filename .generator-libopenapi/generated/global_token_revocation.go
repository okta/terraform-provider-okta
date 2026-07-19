// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// GlobalTokenRevocation represents the GlobalTokenRevocation schema
type GlobalTokenRevocation struct {
	// Authentication method <br> **Note:** Currently, only the `SIGNED_JWT` method is supported
	AuthMethod string `json:"authMethod"`
	// URL of the authorization server's global token revocation endpoint
	Endpoint string `json:"endpoint"`
	// Allow partial support for Universal Logout
	PartialLogout bool `json:"partialLogout,omitempty"`
	// The format of the subject
	SubjectFormat string `json:"subjectFormat"`
}
