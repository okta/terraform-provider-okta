// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UIElement represents the UIElement schema
// Specifies the configuration of an input field on an enrollment form
type UIElement struct {
	// UI Schema element options object
	Options map[string]interface{} `json:"options,omitempty"`
	// Specifies the property bound to the input field. It must follow the format `#/properties/PROPERTY_NAME` where `PROPERTY_NAME` is a variable name for an attribute in `profile editor`.
	Scope string `json:"scope,omitempty"`
	// Specifies the relationship between this input element and `scope`. The `Control` value specifies that this input controls the value represented by `scope`.
	Type string `json:"type,omitempty"`
	// Label name for the UI element
	Label string `json:"label,omitempty"`
}
