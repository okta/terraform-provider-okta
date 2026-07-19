// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// GracePeriod represents the GracePeriod schema
// Represents the Grace Period configuration for the device assurance policy
type GracePeriod struct {
	Expiry interface{} `json:"expiry,omitempty"`
	// Represents the type of Grace Period configured for the device assurance policy
	Type string `json:"type,omitempty"`
}
