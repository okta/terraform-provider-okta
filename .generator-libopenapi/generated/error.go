// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// Error represents the Error schema
type Error struct {
	// An Okta code for this type of error
	ErrorLink string `json:"errorLink,omitempty"`
	// A short description of what caused this error. Sometimes this contains dynamically-generated information about your specific error.
	ErrorSummary string `json:"errorSummary,omitempty"`
	ErrorCauses []interface{} `json:"errorCauses,omitempty"`
	// An Okta code for this type of error
	ErrorCode string `json:"errorCode,omitempty"`
	// A unique identifier for this error. This can be used by Okta Support to help with troubleshooting.
	ErrorId string `json:"errorId,omitempty"`
}
