// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OrgPreferences represents the OrgPreferences schema
type OrgPreferences struct {
	// Indicates if the footer is shown on the End-User Dashboard
	ShowEndUserFooter bool `json:"showEndUserFooter,omitempty"`
	// Specifies link relations (see [Web Linking](https://www.rfc-editor.org/rfc/rfc8288)) available for this object using the [JSON Hypertext Application Language](https://datatracker.ietf.org/doc/html/...
	Links map[string]interface{} `json:"_links,omitempty"`
}
