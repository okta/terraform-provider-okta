// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ProfileMappingSource represents the ProfileMappingSource schema
// The parameter is the source of a profile mapping and is a valid [JSON Schema Draft 4](https://datatracker.ietf.org/doc/html/draft-zyp-json-schema-04) document with the following properties. The dat...
type ProfileMappingSource struct {
	// Unique identifier for the application instance or userType
	ID string `json:"id,omitempty"`
	// Variable name of the application instance or name of the referenced UserType
	Name string `json:"name,omitempty"`
	// Type of user referenced in the mapping
	Type string `json:"type,omitempty"`
	Links interface{} `json:"_links,omitempty"`
}
