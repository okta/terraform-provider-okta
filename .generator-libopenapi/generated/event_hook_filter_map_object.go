// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// EventHookFilterMapObject represents the EventHookFilterMapObject schema
type EventHookFilterMapObject struct {
	// The filtered event type
	Event string `json:"event,omitempty"`
	Condition interface{} `json:"condition,omitempty"`
}
