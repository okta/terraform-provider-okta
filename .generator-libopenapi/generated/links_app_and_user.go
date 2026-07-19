// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// LinksAppAndUser represents the LinksAppAndUser schema
// Specifies link relations (see [Web Linking](https://www.rfc-editor.org/rfc/rfc8288)) available using the [JSON Hypertext Application Language](https://datatracker.ietf.org/doc/html/draft-kelly-json...
type LinksAppAndUser struct {
	App interface{} `json:"app,omitempty"`
	Group interface{} `json:"group,omitempty"`
	User interface{} `json:"user,omitempty"`
}
