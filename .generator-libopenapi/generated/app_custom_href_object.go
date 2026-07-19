// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AppCustomHrefObject represents the AppCustomHrefObject schema
type AppCustomHrefObject struct {
	// Describes allowed HTTP verbs for the `href`
	Hints map[string]interface{} `json:"hints,omitempty"`
	// Link URI
	Href string `json:"href"`
	// Link name
	Title string `json:"title,omitempty"`
	// The media type of the link. If omitted, it is implicitly `application/json`.
	Type string `json:"type,omitempty"`
}
