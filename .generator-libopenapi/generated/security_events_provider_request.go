// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SecurityEventsProviderRequest represents the SecurityEventsProviderRequest schema
// The request schema for creating or updating a security events provider. The `settings` must match one of the schemas.
type SecurityEventsProviderRequest struct {
	// The name of the security events provider instance
	Name string `json:"name"`
	// Information about the security events provider for signal ingestion
	Settings map[string]interface{} `json:"settings"`
	// The app type of the security events provider
	Type string `json:"type"`
}
