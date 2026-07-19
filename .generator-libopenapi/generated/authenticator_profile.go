// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AuthenticatorProfile represents the AuthenticatorProfile schema
// Defines the authenticator specific parameters
type AuthenticatorProfile struct {
	// The phone number for a `call` or `sms` authenticator enrollment.
	PhoneNumber string `json:"phoneNumber"`
}
