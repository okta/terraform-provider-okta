// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UISchemaObject represents the UISchemaObject schema
// Properties of the UI schema
type UISchemaObject struct {
	// Specifies the button label for the `Submit` button at the bottom of the enrollment form
	ButtonLabel string `json:"buttonLabel,omitempty"`
	Elements []interface{} `json:"elements,omitempty"`
	// Specifies the label at the top of the enrollment form under the logo
	Label string `json:"label,omitempty"`
	// Specifies the type of layout
	Type string `json:"type,omitempty"`
}
