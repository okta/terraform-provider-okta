// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SelfServicePasswordResetAction represents the SelfServicePasswordResetAction schema
// Enables or disables users to reset their own password and defines the authenticators and constraints needed to complete the reset
type SelfServicePasswordResetAction struct {
	Access interface{} `json:"access,omitempty"`
	Requirement interface{} `json:"requirement,omitempty"`
	// <x-lifecycle class="oie"></x-lifecycle> The type of rule action
	Type string `json:"type,omitempty"`
}
