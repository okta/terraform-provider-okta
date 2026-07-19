// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// LogUserBehavior represents the LogUserBehavior schema
// The result of the user behavior detection models associated with the event
type LogUserBehavior struct {
	// The unique identifier of the user behavior detection model
	ID string `json:"id,omitempty"`
	// The name of the user behavior detection model [configured by admins](https://developer.okta.com/docs/api/openapi/okta-management/management/tag/Behavior/)
	Name string `json:"name,omitempty"`
	// The result of the user behavior analysis
	Result string `json:"result,omitempty"`
}
