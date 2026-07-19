// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// LogDebugContext represents the LogDebugContext schema
// For some kinds of events (for example, OLM provisioning, sign-in request, second factor SMS, and so on), the fields that are provided in other response objects aren't sufficient to adequately descr...
type LogDebugContext struct {
	// A dynamic field that contains miscellaneous information that is dependent on the event type.
	DebugData map[string]interface{} `json:"debugData,omitempty"`
}
