// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// Error409 represents the Error409 schema
// Conflict error object
type Error409 struct {
	// Another request has already been received for the settings for this email template
	ErrorCauses []interface{} `json:"errorCauses,omitempty"`
	// E0000254
	ErrorCode string `json:"errorCode,omitempty"`
	// sampleH3iLB6bpBcbnV9E09Fy
	ErrorId string `json:"errorId,omitempty"`
	// E0000254
	ErrorLink string `json:"errorLink,omitempty"`
	// Another request has already been received for the settings for this email template
	ErrorSummary string `json:"errorSummary,omitempty"`
}
