// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// GroupSchemaCustom represents the GroupSchemaCustom schema
// All custom profile properties are defined in a profile subschema with the resolution scope `#custom`
type GroupSchemaCustom struct {
	// The subschema name
	ID string `json:"id,omitempty"`
	// The `#custom` object properties
	Properties map[string]interface{} `json:"properties,omitempty"`
	// A collection indicating required property names
	Required []string `json:"required,omitempty"`
	// The object type
	Type string `json:"type,omitempty"`
}
