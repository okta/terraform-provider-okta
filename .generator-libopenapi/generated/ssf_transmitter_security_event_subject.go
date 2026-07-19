// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SsfTransmitterSecurityEventSubject represents the SsfTransmitterSecurityEventSubject schema
// The event subject
type SsfTransmitterSecurityEventSubject struct {
	// The format of the subject
	Format string `json:"format,omitempty"`
	// An identifier of the actor
	Iss string `json:"iss,omitempty"`
	// An identifier for the subject that was acted on
	Sub string `json:"sub,omitempty"`
}
