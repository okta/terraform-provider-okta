// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// WellKnownOrgMetadata represents the WellKnownOrgMetadata schema
type WellKnownOrgMetadata struct {
	// Org unique identifier
	ID string `json:"id,omitempty"`
	Pipeline interface{} `json:"pipeline,omitempty"`
	// Specifies link relations (see [Web Linking](https://www.rfc-editor.org/rfc/rfc8288)) available for this object using the [JSON Hypertext Application Language](https://datatracker.ietf.org/doc/html/...
	Links map[string]interface{} `json:"_links,omitempty"`
}
