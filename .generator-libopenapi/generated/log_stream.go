// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// LogStream represents the LogStream schema
type LogStream struct {
	// Timestamp when the log stream object was last updated
	LastUpdated *time.Time `json:"lastUpdated"`
	Name interface{} `json:"name"`
	// Lifecycle status of the log stream object
	Status string `json:"status"`
	Type interface{} `json:"type"`
	Links interface{} `json:"_links"`
	// Timestamp when the log stream object was created
	Created *time.Time `json:"created"`
	// Unique identifier for the log stream
	ID string `json:"id"`
}
