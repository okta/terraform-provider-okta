// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// GroupSchema represents the GroupSchema schema
type GroupSchema struct {
	// Name of the schema
	Name string `json:"name,omitempty"`
	// Type of [root schema](https://tools.ietf.org/html/draft-zyp-json-schema-04#section-3.4)
	Type string `json:"type,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	// URI of group schema
	ID string `json:"id,omitempty"`
	// Group object properties
	Properties interface{} `json:"properties,omitempty"`
	// User-defined display name for the schema
	Title string `json:"title,omitempty"`
	// JSON schema version identifier
	$schema string `json:"$schema,omitempty"`
	// Timestamp when the schema was created
	Created string `json:"created,omitempty"`
	Definitions interface{} `json:"definitions,omitempty"`
	// Description for the schema
	Description string `json:"description,omitempty"`
	// Timestamp when the schema was last updated
	LastUpdated string `json:"lastUpdated,omitempty"`
}
