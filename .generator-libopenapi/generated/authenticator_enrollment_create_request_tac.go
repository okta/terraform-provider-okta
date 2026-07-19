// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AuthenticatorEnrollmentCreateRequestTac represents the AuthenticatorEnrollmentCreateRequestTac schema
type AuthenticatorEnrollmentCreateRequestTac struct {
	// Unique identifier of the TAC authenticator
	AuthenticatorId string `json:"authenticatorId"`
	Profile interface{} `json:"profile,omitempty"`
}
