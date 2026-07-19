// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// IntegrationCapability represents the IntegrationCapability schema
type IntegrationCapability struct {
	// Name of the integration capability
	Name string `json:"name,omitempty"`
	// The date when the integration capability was released
	ReleaseDate *time.Time `json:"releaseDate,omitempty"`
	// Status of the integration capability
	Status string `json:"status,omitempty"`
	// Description of the integration capability
	Description string `json:"description,omitempty"`
	// URL to the help documentation for the integration capability
	HelpUrl string `json:"helpUrl,omitempty"`
	// Unique identifier for the integration capability
	ID string `json:"id,omitempty"`
}
