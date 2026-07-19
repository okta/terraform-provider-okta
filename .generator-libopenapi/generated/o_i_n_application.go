// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OINApplication represents the OINApplication schema
type OINApplication struct {
	Label interface{} `json:"label,omitempty"`
	// The key name for the OIN app definition
	Name string `json:"name,omitempty"`
	Accessibility interface{} `json:"accessibility,omitempty"`
	Credentials interface{} `json:"credentials,omitempty"`
	Licensing interface{} `json:"licensing,omitempty"`
	// Contains any valid JSON schema for specifying properties that can be referenced from a request (only available to OAuth 2.0 client apps)
	Profile map[string]interface{} `json:"profile,omitempty"`
	// Authentication mode for the app
	SignOnMode string `json:"signOnMode,omitempty"`
	Status interface{} `json:"status,omitempty"`
	Visibility interface{} `json:"visibility,omitempty"`
}
