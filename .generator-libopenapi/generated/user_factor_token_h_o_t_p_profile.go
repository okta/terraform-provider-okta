// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserFactorTokenHOTPProfile represents the UserFactorTokenHOTPProfile schema
type UserFactorTokenHOTPProfile struct {
	// Unique secret key used to generate the OTP
	SharedSecret string `json:"sharedSecret,omitempty"`
}
