// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SecurityEventSubject represents the SecurityEventSubject schema
// The event subjects
type SecurityEventSubject struct {
	// The device involved with the event
	Device map[string]interface{} `json:"device,omitempty"`
	// The user involved with the event
	User map[string]interface{} `json:"user,omitempty"`
}
