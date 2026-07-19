// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// OSAccount represents the OSAccount schema
type OSAccount struct {
	// Unique identifier of the device this OS account belongs to
	DeviceId string `json:"deviceId"`
	// Unique identifier for the OS account
	ID string `json:"id"`
	// Timestamp when the OS account was last updated
	LastUpdated *time.Time `json:"lastUpdated"`
	Platform interface{} `json:"platform"`
	Links interface{} `json:"_links"`
	// Timestamp when the OS account was created
	Created *time.Time `json:"created"`
}
