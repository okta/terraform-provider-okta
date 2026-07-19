// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// FulfillmentRequest represents the FulfillmentRequest schema
// Fulfillment request
type FulfillmentRequest struct {
	FulfillmentData interface{} `json:"fulfillmentData,omitempty"`
	// Name of the fulfillment provider for the WebAuthn preregistration factor
	FulfillmentProvider string `json:"fulfillmentProvider,omitempty"`
	// ID of an existing Okta user
	UserId string `json:"userId,omitempty"`
}
