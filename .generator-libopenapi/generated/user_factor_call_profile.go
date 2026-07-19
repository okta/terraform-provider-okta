// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserFactorCallProfile represents the UserFactorCallProfile schema
type UserFactorCallProfile struct {
	// Extension of the associated `phoneNumber`
	PhoneExtension string `json:"phoneExtension,omitempty"`
	// Phone number of the factor. Format phone numbers to use the [E.164 standard](https://www.itu.int/rec/T-REC-E.164/).
	PhoneNumber string `json:"phoneNumber,omitempty"`
}
