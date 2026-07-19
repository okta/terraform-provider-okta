// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AuthenticatorEnrollmentCreateRequest represents the AuthenticatorEnrollmentCreateRequest schema
type AuthenticatorEnrollmentCreateRequest struct {
	// Unique identifier of the `phone` authenticator
	AuthenticatorId string `json:"authenticatorId"`
	Profile interface{} `json:"profile"`
}
