// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// orgTechnicalContactType represents the orgTechnicalContactType schema
// Org technical contact
type orgTechnicalContactType struct {
	ContactType interface{} `json:"contactType,omitempty"`
	// Specifies link relations (see [Web Linking](https://www.rfc-editor.org/rfc/rfc8288)) available for the org technical Contact Type object using the [JSON Hypertext Application Language](https://data...
	Links map[string]interface{} `json:"_links,omitempty"`
}
