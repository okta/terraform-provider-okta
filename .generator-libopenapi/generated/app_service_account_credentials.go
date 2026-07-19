// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AppServiceAccountCredentials represents the AppServiceAccountCredentials schema
// Credentials for a SaaS app account
type AppServiceAccountCredentials struct {
	// The password associated with the service account
	Password string `json:"password,omitempty"`
	// The username associated with the service account
	Username string `json:"username"`
}
