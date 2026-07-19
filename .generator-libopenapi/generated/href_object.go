// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// HrefObject represents the HrefObject schema
type HrefObject struct {
	Hints interface{} `json:"hints,omitempty"`
	// Link URI
	Href string `json:"href"`
	// Link name
	Name string `json:"name,omitempty"`
	// Indicates whether the link object's `href` property is a URI template.
	Templated bool `json:"templated,omitempty"`
	// The media type of the link. If omitted, it is implicitly `application/json`.
	Type string `json:"type,omitempty"`
}
