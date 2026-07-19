// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// GroupSchemaBase represents the GroupSchemaBase schema
type GroupSchemaBase struct {
	// The subschema name
	ID string `json:"id,omitempty"`
	// The `#base` object properties
	Properties interface{} `json:"properties,omitempty"`
	// A collection indicating required property names
	Required []string `json:"required,omitempty"`
	// The object type
	Type string `json:"type,omitempty"`
}
