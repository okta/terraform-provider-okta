// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OrgContactUser represents the OrgContactUser schema
type OrgContactUser struct {
	// Contact user ID
	UserId string `json:"userId,omitempty"`
	// Specifies link relations (see [Web Linking](https://www.rfc-editor.org/rfc/rfc8288)) available for the contact type user object using the [JSON Hypertext Application Language](https://datatracker.i...
	Links map[string]interface{} `json:"_links,omitempty"`
}
