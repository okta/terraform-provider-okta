// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AuthenticatorMethodConstraint represents the AuthenticatorMethodConstraint schema
// Limits the authenticators that can be used for a given method. Currently, only the `otp` method supports constraints, and Google authenticator (key : 'google_otp') is the only allowed authenticator.
type AuthenticatorMethodConstraint struct {
	AllowedAuthenticators []interface{} `json:"allowedAuthenticators,omitempty"`
	Method interface{} `json:"method,omitempty"`
}
