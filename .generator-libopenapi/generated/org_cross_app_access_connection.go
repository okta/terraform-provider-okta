// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// OrgCrossAppAccessConnection represents the OrgCrossAppAccessConnection schema
// Connection object for Cross App Access connections
type OrgCrossAppAccessConnection struct {
	// Unique identifier for the connection
	ID string `json:"id,omitempty"`
	// The ISO 8601 formatted date and time when the connection was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// ID of the requesting app instance
	RequestingAppInstanceId string `json:"requestingAppInstanceId,omitempty"`
	// ID of the resource app instance
	ResourceAppInstanceId string `json:"resourceAppInstanceId,omitempty"`
	// Indicates if the Cross App Access connection is active or inactive
	Status string `json:"status,omitempty"`
	// The ISO 8601 formatted date and time when the connection was created
	Created *time.Time `json:"created,omitempty"`
}
