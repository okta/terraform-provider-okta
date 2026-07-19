// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ListProfileMappings represents the ListProfileMappings schema
// A collection of the profile mappings that include a subset of the profile mapping object's properties. The profile mapping object describes a mapping between an Okta user's and an app user's proper...
type ListProfileMappings struct {
	Source interface{} `json:"source,omitempty"`
	Target interface{} `json:"target,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	// Unique identifier for profile mapping
	ID string `json:"id,omitempty"`
}
