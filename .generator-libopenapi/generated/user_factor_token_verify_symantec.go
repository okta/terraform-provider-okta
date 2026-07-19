// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserFactorTokenVerifySymantec represents the UserFactorTokenVerifySymantec schema
type UserFactorTokenVerifySymantec struct {
	// OTP for the current time window
	PassCode string `json:"passCode,omitempty"`
	// OTP for the next time window
	NextPassCode int `json:"nextPassCode,omitempty"`
}
