// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// RiskEvent represents the RiskEvent schema
type RiskEvent struct {
	// Timestamp at which the event expires (expressed as a UTC time zone using ISO 8601 format: yyyy-MM-dd`T`HH:mm:ss.SSS`Z`). If this optional field isn't included, Okta automatically expires the event ...
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	// List of risk event subjects
	Subjects []interface{} `json:"subjects"`
	// Timestamp of when the event is produced (expressed as a UTC time zone using ISO 8601 format: yyyy-MM-dd`T`HH:mm:ss.SSS`Z`)
	Timestamp *time.Time `json:"timestamp,omitempty"`
}
