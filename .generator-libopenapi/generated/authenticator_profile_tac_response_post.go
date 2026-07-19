// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// AuthenticatorProfileTacResponsePost represents the AuthenticatorProfileTacResponsePost schema
// Defines the authenticator specific parameters
type AuthenticatorProfileTacResponsePost struct {
	// The time when the TAC enrollment expires in the UTC timezone
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	// Determines whether an enrollment can be used more than once
	MultiUse bool `json:"multiUse,omitempty"`
	// A temporary access code used for authentication. It can be used one or more times and is valid for a defined period specified by the `ttl` property. The `tac` is returned in the response when the e...
	Tac string `json:"tac,omitempty"`
}
