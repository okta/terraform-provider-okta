// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AuthenticatorProfileTacRequest represents the AuthenticatorProfileTacRequest schema
// Defines the authenticator specific parameters
type AuthenticatorProfileTacRequest struct {
	// Determines whether the enrollment can be used more than once. To enable multi-use, the org-level authenticator’s configuration must allow multi-use.
	MultiUse bool `json:"multiUse,omitempty"`
	// Time-to-live (TTL) in minutes.  Specifies how long the TAC enrollment is valid after it's created and activated. The configured value must be between 10 minutes (`10`) and 10 days (`14400`), inclus...
	Ttl string `json:"ttl,omitempty"`
}
