// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// Feature represents the Feature schema
// Specifies feature release cycle information
type Feature struct {
	Type interface{} `json:"type,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	// Brief description of the feature and what it provides
	Description string `json:"description,omitempty"`
	// Unique identifier for this feature
	ID string `json:"id,omitempty"`
	// Name of the feature
	Name string `json:"name,omitempty"`
	Stage interface{} `json:"stage,omitempty"`
	Status interface{} `json:"status,omitempty"`
}
