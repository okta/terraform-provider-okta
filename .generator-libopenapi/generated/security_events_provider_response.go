// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SecurityEventsProviderResponse represents the SecurityEventsProviderResponse schema
// The security events provider response
type SecurityEventsProviderResponse struct {
	// The name of the security events provider instance
	Name string `json:"name,omitempty"`
	// Information about the security events provider for signal ingestion
	Settings interface{} `json:"settings,omitempty"`
	// Indicates whether the security events provider is active or not
	Status string `json:"status,omitempty"`
	// The app type of the security events provider
	Type string `json:"type,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	// The unique identifier of this instance
	ID string `json:"id,omitempty"`
}
