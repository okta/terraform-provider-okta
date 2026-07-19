// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ProfileMappingTarget represents the ProfileMappingTarget schema
// The parameter is the target of a profile mapping and is a valid [JSON Schema Draft 4](https://datatracker.ietf.org/doc/html/draft-zyp-json-schema-04) document with the following properties. The dat...
type ProfileMappingTarget struct {
	// Unique identifier for the application instance or UserType
	ID string `json:"id,omitempty"`
	// Variable name of the application instance or name of the referenced userType
	Name string `json:"name,omitempty"`
	// Type of user referenced in the mapping
	Type string `json:"type,omitempty"`
	Links interface{} `json:"_links,omitempty"`
}
