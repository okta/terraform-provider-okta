// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// InlineHookBasePayload represents the InlineHookBasePayload schema
type InlineHookBasePayload struct {
	// The time the inline hook request was sent
	EventTime string `json:"eventTime,omitempty"`
	// The inline hook version
	EventTypeVersion string `json:"eventTypeVersion,omitempty"`
	// The inline hook cloud version
	CloudEventVersion string `json:"cloudEventVersion,omitempty"`
	// The inline hook request header content
	ContentType string `json:"contentType,omitempty"`
	// The individual inline hook request ID
	EventId string `json:"eventId,omitempty"`
}
