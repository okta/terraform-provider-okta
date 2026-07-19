// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// PrivilegedResource represents the PrivilegedResource schema
// Base class for PrivilegedResourceRequest and PrivilegedResourceResponse
type PrivilegedResource struct {
	CredentialSyncInfo interface{} `json:"credentialSyncInfo,omitempty"`
	// ID of the privileged resource
	ID string `json:"id,omitempty"`
	// Timestamp when the object was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	ResourceType interface{} `json:"resourceType"`
	Status interface{} `json:"status,omitempty"`
	// Timestamp when the object was created
	Created *time.Time `json:"created,omitempty"`
}
