// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserFactorSupported represents the UserFactorSupported schema
type UserFactorSupported struct {
	// Indicates if the factor is required for the specified user
	Enrollment string `json:"enrollment,omitempty"`
	FactorType interface{} `json:"factorType,omitempty"`
	Provider interface{} `json:"provider,omitempty"`
	Status interface{} `json:"status,omitempty"`
	// Name of the factor vendor. This is usually the same as the provider except for On-Prem MFA, which depends on admin settings.
	VendorName string `json:"vendorName,omitempty"`
	// Embedded resources related to the factor
	Embedded map[string]interface{} `json:"_embedded,omitempty"`
	Links interface{} `json:"_links,omitempty"`
}
