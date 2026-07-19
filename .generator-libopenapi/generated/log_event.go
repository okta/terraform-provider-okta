// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// LogEvent represents the LogEvent schema
type LogEvent struct {
	// The entity that an actor performs an action on. Targets can be anything, such as an app user, a sign-in token, or anything else.  > **Note:** When searching the target array, search for a given `ty...
	Target []interface{} `json:"target,omitempty"`
	Transaction interface{} `json:"transaction,omitempty"`
	Actor interface{} `json:"actor,omitempty"`
	DebugContext interface{} `json:"debugContext,omitempty"`
	// Timestamp when the event is published
	Published *time.Time `json:"published,omitempty"`
	AuthenticationContext interface{} `json:"authenticationContext,omitempty"`
	// The display message for an event
	DisplayMessage string `json:"displayMessage,omitempty"`
	// Versioning indicator
	Version string `json:"version,omitempty"`
	Client interface{} `json:"client,omitempty"`
	// The published event type. Event instances are categorized by action in the event type attribute. This attribute is key to navigating the System Log through expression filters. See [Event Types cata...
	EventType string `json:"eventType,omitempty"`
	// Associated Events API Action `objectType` attribute value
	LegacyEventType string `json:"legacyEventType,omitempty"`
	Outcome interface{} `json:"outcome,omitempty"`
	Severity interface{} `json:"severity,omitempty"`
	// Unique identifier for an individual event
	Uuid string `json:"uuid,omitempty"`
	Request interface{} `json:"request,omitempty"`
	SecurityContext interface{} `json:"securityContext,omitempty"`
}
