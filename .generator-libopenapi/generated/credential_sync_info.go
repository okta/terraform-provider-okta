// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// CredentialSyncInfo represents the CredentialSyncInfo schema
type CredentialSyncInfo struct {
	// The error code for the type of error
	ErrorCode string `json:"errorCode,omitempty"`
	// A short description of the error
	ErrorReason string `json:"errorReason,omitempty"`
	// The version ID of the password secret from the OPA vault.
	SecretVersionId string `json:"secretVersionId,omitempty"`
	SyncState interface{} `json:"syncState,omitempty"`
	// Timestamp when the credential was changed
	SyncTime *time.Time `json:"syncTime,omitempty"`
}
