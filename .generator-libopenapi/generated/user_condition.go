// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserCondition represents the UserCondition schema
// Specifies a set of users to be included or excluded
type UserCondition struct {
	// Users to be included
	Include []string `json:"include,omitempty"`
	// Users to be excluded
	Exclude []string `json:"exclude,omitempty"`
}
