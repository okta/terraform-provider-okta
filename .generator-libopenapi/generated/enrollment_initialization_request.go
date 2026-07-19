// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// EnrollmentInitializationRequest represents the EnrollmentInitializationRequest schema
// Enrollment initialization request
type EnrollmentInitializationRequest struct {
	// Name of the fulfillment provider for the WebAuthn preregistration factor
	FulfillmentProvider string `json:"fulfillmentProvider,omitempty"`
	// ID of an existing Okta user
	UserId string `json:"userId,omitempty"`
	// Transport public key in JWK (JSON Web Key) format used to encrypt fulfillment requests to Yubico
	YubicoTransportKeyJWK interface{} `json:"yubicoTransportKeyJWK,omitempty"`
	// List of relying party hostnames to register on the YubiKey
	EnrollmentRpIds []string `json:"enrollmentRpIds,omitempty"`
}
