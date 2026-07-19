// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SupportedMethods represents the SupportedMethods schema
// The supported methods of an authenticator
type SupportedMethods struct {
	// The type of authenticator method
	Type string `json:"type,omitempty"`
	Settings map[string]interface{} `json:"settings,omitempty"`
	Status interface{} `json:"status,omitempty"`
}
