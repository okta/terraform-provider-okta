// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// WebAuthnPreregistrationFactor represents the WebAuthnPreregistrationFactor schema
// User factor variant used for WebAuthn preregistration factors
type WebAuthnPreregistrationFactor struct {
	Provider interface{} `json:"provider,omitempty"`
	// Name of the factor vendor. This is usually the same as the provider.
	VendorName string `json:"vendorName,omitempty"`
	// Timestamp indicating when the factor was enrolled
	Created *time.Time `json:"created,omitempty"`
	// ID of the factor
	ID string `json:"id,omitempty"`
	Status interface{} `json:"status,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	FactorType interface{} `json:"factorType,omitempty"`
	// Timestamp indicating when the factor was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// Specific attributes related to the factor
	Profile map[string]interface{} `json:"profile,omitempty"`
}
