// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// Conditions represents the Conditions schema
// Conditions of applying realm assignment
type Conditions struct {
	// ID of the profile source
	ProfileSourceId string `json:"profileSourceId,omitempty"`
	Expression interface{} `json:"expression,omitempty"`
}
