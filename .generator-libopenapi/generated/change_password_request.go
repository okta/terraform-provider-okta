// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ChangePasswordRequest represents the ChangePasswordRequest schema
type ChangePasswordRequest struct {
	// When set to `true`, revokes all user sessions, except for the current session
	RevokeSessions bool `json:"revokeSessions,omitempty"`
	NewPassword interface{} `json:"newPassword,omitempty"`
	OldPassword interface{} `json:"oldPassword,omitempty"`
}
