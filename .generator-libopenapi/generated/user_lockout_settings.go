// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserLockoutSettings represents the UserLockoutSettings schema
type UserLockoutSettings struct {
	// Prevents brute-force lockout from unknown devices for the password authenticator.
	PreventBruteForceLockoutFromUnknownDevices bool `json:"preventBruteForceLockoutFromUnknownDevices,omitempty"`
}
