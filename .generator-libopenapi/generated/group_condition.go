// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// GroupCondition represents the GroupCondition schema
// Specifies a set of groups whose users are to be included or excluded
type GroupCondition struct {
	// Groups to be excluded
	Exclude []string `json:"exclude,omitempty"`
	// Groups to be included
	Include []string `json:"include,omitempty"`
}
