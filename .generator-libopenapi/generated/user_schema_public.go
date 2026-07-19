// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserSchemaPublic represents the UserSchemaPublic schema
// All custom profile properties are defined in a profile subschema with the resolution scope `#custom`.  > **Notes:** > * When you refer to custom profile attributes that differ only by case, name co...
type UserSchemaPublic struct {
	// The subschema name
	ID string `json:"id,omitempty"`
	// The `#custom` object properties
	Properties map[string]interface{} `json:"properties,omitempty"`
	// A collection indicating required property names
	Required []string `json:"required,omitempty"`
	// The object type
	Type string `json:"type,omitempty"`
}
