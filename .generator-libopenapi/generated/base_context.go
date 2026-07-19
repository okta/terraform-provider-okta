// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// BaseContext represents the BaseContext schema
// This object contains a number of sub-objects, each of which provide some type of contextual information.
type BaseContext struct {
	Request interface{} `json:"request,omitempty"`
	// Details of the user session
	Session map[string]interface{} `json:"session,omitempty"`
	// Identifies the Okta user that the token was generated to authenticate and provides details of their Okta user profile
	User map[string]interface{} `json:"user,omitempty"`
}
