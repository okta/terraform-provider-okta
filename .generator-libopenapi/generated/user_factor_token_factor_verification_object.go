// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserFactorTokenFactorVerificationObject represents the UserFactorTokenFactorVerificationObject schema
type UserFactorTokenFactorVerificationObject struct {
	// OTP for the next time window
	NextPassCode string `json:"nextPassCode,omitempty"`
	// OTP for the current time window
	PassCode string `json:"passCode,omitempty"`
}
