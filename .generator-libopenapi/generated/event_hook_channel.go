// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// EventHookChannel represents the EventHookChannel schema
type EventHookChannel struct {
	// Version of the channel. Currently the only supported version is `1.0.0`.
	Version string `json:"version"`
	Config interface{} `json:"config"`
	Type interface{} `json:"type"`
}
