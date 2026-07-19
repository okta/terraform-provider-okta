// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// LogStreamLinksSelfAndLifecycle represents the LogStreamLinksSelfAndLifecycle schema
// Specifies link relations (see [Web Linking](https://www.rfc-editor.org/rfc/rfc8288)) available for the current status of an application using the [JSON Hypertext Application Language](https://datat...
type LogStreamLinksSelfAndLifecycle struct {
	Activate interface{} `json:"activate,omitempty"`
	Deactivate interface{} `json:"deactivate,omitempty"`
	Self interface{} `json:"self"`
}
