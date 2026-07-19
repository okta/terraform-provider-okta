// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// LogActor represents the LogActor schema
// Describes the user, app, client, or other entity (actor) who performs an action on a target. The actor is dependent on the action that is performed. All events have actors.
type LogActor struct {
	// Further details about the actor
	DetailEntry map[string]interface{} `json:"detailEntry,omitempty"`
	// Display name of the actor
	DisplayName string `json:"displayName,omitempty"`
	// ID of the actor
	ID string `json:"id,omitempty"`
	// Type of actor
	Type string `json:"type,omitempty"`
	// Alternative ID of the actor
	AlternateId string `json:"alternateId,omitempty"`
}
