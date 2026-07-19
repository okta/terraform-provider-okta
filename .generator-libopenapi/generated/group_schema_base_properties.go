// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// GroupSchemaBaseProperties represents the GroupSchemaBaseProperties schema
// All Okta-defined profile properties are defined in a profile subschema with the resolution scope `#base`. These properties can't be removed or edited, regardless of any attempt to do so.
type GroupSchemaBaseProperties struct {
	// Unique identifier for the group
	Name interface{} `json:"name,omitempty"`
	// Human readable description of the group
	Description interface{} `json:"description,omitempty"`
}
