// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// DeviceDisplayName represents the DeviceDisplayName schema
// Display name of the device
type DeviceDisplayName struct {
	// Indicates whether the associated value is Personal Identifiable Information (PII) and requires masking
	Sensitive bool `json:"sensitive,omitempty"`
	// Display name of the device
	Value string `json:"value,omitempty"`
}
