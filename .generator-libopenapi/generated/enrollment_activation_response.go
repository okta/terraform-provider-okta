// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// EnrollmentActivationResponse represents the EnrollmentActivationResponse schema
// Enrollment initialization response
type EnrollmentActivationResponse struct {
	// ID of an existing Okta user
	UserId string `json:"userId,omitempty"`
	// List of IDs for preregistered WebAuthn factors in Okta
	AuthenticatorEnrollmentIds []string `json:"authenticatorEnrollmentIds,omitempty"`
	// Name of the fulfillment provider for the WebAuthn preregistration factor
	FulfillmentProvider string `json:"fulfillmentProvider,omitempty"`
}
