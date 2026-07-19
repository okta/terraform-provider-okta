// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// EnrollmentInitializationResponse represents the EnrollmentInitializationResponse schema
// Yubico transport key in the form of a JSON Web Token (JWK), used to encrypt our fulfillment request to Yubico. The currently agreed protocol uses P-384.
type EnrollmentInitializationResponse struct {
	// ID of an existing Okta user
	UserId string `json:"userId,omitempty"`
	// List of credential requests for the fulfillment provider
	CredRequests []interface{} `json:"credRequests,omitempty"`
	// Name of the fulfillment provider for the WebAuthn preregistration factor
	FulfillmentProvider string `json:"fulfillmentProvider,omitempty"`
	// Encrypted JWE of PIN request for the fulfillment provider
	PinRequestJwe string `json:"pinRequestJwe,omitempty"`
}
