// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OAuth2RefreshTokenScope represents the OAuth2RefreshTokenScope schema
type OAuth2RefreshTokenScope struct {
	// Name of the end user displayed in a consent dialog
	DisplayName string `json:"displayName,omitempty"`
	// Scope object ID
	ID string `json:"id,omitempty"`
	// Scope name
	Name string `json:"name,omitempty"`
	// Specifies link relations (see [Web Linking](https://www.rfc-editor.org/rfc/rfc8288)) available for the current status of an application using the [JSON Hypertext Application Language](https://datat...
	Links map[string]interface{} `json:"_links,omitempty"`
	// Description of the Scope
	Description string `json:"description,omitempty"`
}
