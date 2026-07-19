// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PinRequest represents the PinRequest schema
// Pin request
type PinRequest struct {
	// ID for a WebAuthn preregistration factor in Okta
	AuthenticatorEnrollmentId string `json:"authenticatorEnrollmentId,omitempty"`
	// Name of the fulfillment provider for the WebAuthn preregistration factor
	FulfillmentProvider string `json:"fulfillmentProvider,omitempty"`
	// ID of an existing Okta user
	UserId string `json:"userId,omitempty"`
}
