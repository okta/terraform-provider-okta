// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserSchema represents the UserSchema schema
type UserSchema struct {
	// User Object Properties
	Properties interface{} `json:"properties,omitempty"`
	// Type of [root schema](https://tools.ietf.org/html/draft-zyp-json-schema-04#section-3.4)
	Type string `json:"type,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	// JSON schema version identifier
	$schema string `json:"$schema,omitempty"`
	// Timestamp when the schema was created
	Created string `json:"created,omitempty"`
	// URI of user schema
	ID string `json:"id,omitempty"`
	// Timestamp when the schema was last updated
	LastUpdated string `json:"lastUpdated,omitempty"`
	// User-defined display name for the schema
	Title string `json:"title,omitempty"`
	// User profile subschemas  The profile object for a user is defined by a composite schema of base and custom properties using a JSON path to reference subschemas. The `#base` properties are defined a...
	Definitions interface{} `json:"definitions,omitempty"`
	// Name of the schema
	Name string `json:"name,omitempty"`
}
