// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// OperationResponse represents the OperationResponse schema
type OperationResponse struct {
	// Timestamp of when the operation started
	Started *time.Time `json:"started,omitempty"`
	// The status of the asynchronous operation
	Status string `json:"status"`
	// The operation type
	Type string `json:"type"`
	// Timestamp of when the operation completed
	Completed *time.Time `json:"completed,omitempty"`
	// Timestamp of when the operation was created
	Created *time.Time `json:"created"`
	// ID of the asynchronous operation
	ID string `json:"id"`
}
