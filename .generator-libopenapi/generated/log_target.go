// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// LogTarget represents the LogTarget schema
type LogTarget struct {
	// The ID of the target
	ID string `json:"id,omitempty"`
	// The type of target
	Type string `json:"type,omitempty"`
	// The alternate ID of the target
	AlternateId string `json:"alternateId,omitempty"`
	// Details on the target's changes. Not all event types support the `changeDetails` property, and not all `target` objects contain the `changeDetails` property.  > **Note:** You can't run queries on `...
	ChangeDetails map[string]interface{} `json:"changeDetails,omitempty"`
	// Further details on the target
	DetailEntry map[string]interface{} `json:"detailEntry,omitempty"`
	// The display name of the target
	DisplayName string `json:"displayName,omitempty"`
}
