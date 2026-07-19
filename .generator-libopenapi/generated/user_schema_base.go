// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserSchemaBase represents the UserSchemaBase schema
// All Okta-defined profile properties are defined in a profile subschema with the resolution scope `#base`. You can't modify these properties, except to update permissions, to change the nullability ...
type UserSchemaBase struct {
	// The subschema name
	ID string `json:"id,omitempty"`
	// The `#base` object properties
	Properties interface{} `json:"properties,omitempty"`
	// A collection indicating required property names
	Required []string `json:"required,omitempty"`
	// The object type
	Type string `json:"type,omitempty"`
}
