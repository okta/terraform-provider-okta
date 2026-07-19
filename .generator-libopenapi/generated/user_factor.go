// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// UserFactor represents the UserFactor schema
type UserFactor struct {
	Links interface{} `json:"_links,omitempty"`
	// Specific attributes related to the factor
	Profile map[string]interface{} `json:"profile,omitempty"`
	// Provider for the factor. Each provider can support a subset of factor types.
	Provider string `json:"provider,omitempty"`
	// Name of the factor vendor. This is usually the same as the provider except for On-Prem MFA, which depends on admin settings.
	VendorName string `json:"vendorName,omitempty"`
	// Timestamp when the factor was enrolled
	Created *time.Time `json:"created,omitempty"`
	FactorType interface{} `json:"factorType,omitempty"`
	// ID of the factor
	ID string `json:"id,omitempty"`
	// Timestamp when the factor was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	Status interface{} `json:"status,omitempty"`
	Embedded map[string]interface{} `json:"_embedded,omitempty"`
}
