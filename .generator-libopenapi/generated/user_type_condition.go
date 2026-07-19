// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserTypeCondition represents the UserTypeCondition schema
// <x-lifecycle class="oie"></x-lifecycle> Specifies which user types to include and/or exclude
type UserTypeCondition struct {
	// The user types to exclude
	Exclude []string `json:"exclude"`
	// The user types to include
	Include []string `json:"include"`
}
