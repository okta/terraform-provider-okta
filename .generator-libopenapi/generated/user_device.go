// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// UserDevice represents the UserDevice schema
type UserDevice struct {
	// Timestamp when the device was created
	Created *time.Time `json:"created,omitempty"`
	Device map[string]interface{} `json:"device,omitempty"`
	// Unique key for the user device link
	DeviceUserId string `json:"deviceUserId,omitempty"`
}
