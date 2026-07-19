// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SsprRequirement represents the SsprRequirement schema
// <x-lifecycle class="oie"></x-lifecycle> Describes the initial and secondary authenticator requirements a user needs to reset their password
type SsprRequirement struct {
	// Determines which authentication requirements a user needs to perform self-service operations. `AUTH_POLICY` defers conditions and authentication requirements to the [Okta account management policy]...
	AccessControl string `json:"accessControl,omitempty"`
	Primary interface{} `json:"primary,omitempty"`
	StepUp interface{} `json:"stepUp,omitempty"`
}
