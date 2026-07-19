// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// EnrollmentActivationRequest represents the EnrollmentActivationRequest schema
// Enrollment Initialization Request
type EnrollmentActivationRequest struct {
	// Serial number of the YubiKey
	Serial string `json:"serial,omitempty"`
	// ID of an existing Okta user
	UserId string `json:"userId,omitempty"`
	// Firmware version of the YubiKey
	Version string `json:"version,omitempty"`
	// List of usable signing keys from Yubico (in JSON Web Key Sets (JWKS) format). The signing keys are used to verify the JSON Web Signature (JWS) inside the JWE.
	YubicoSigningJwks []interface{} `json:"yubicoSigningJwks,omitempty"`
	// List of credential responses from the fulfillment provider
	CredResponses []interface{} `json:"credResponses,omitempty"`
	// Name of the fulfillment provider for the WebAuthn preregistration factor
	FulfillmentProvider string `json:"fulfillmentProvider,omitempty"`
	// Encrypted JWE of the PIN response from the fulfillment provider
	PinResponseJwe string `json:"pinResponseJwe,omitempty"`
}
