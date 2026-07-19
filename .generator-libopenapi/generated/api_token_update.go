// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// ApiTokenUpdate represents the ApiTokenUpdate schema
// An API Token Update Object for an Okta user. This token is NOT scoped any further and can be used for any API that the user has permissions to call.
type ApiTokenUpdate struct {
	// The creation date of the API Token
	Created *time.Time `json:"created,omitempty"`
	// The name associated with the API Token
	Name string `json:"name,omitempty"`
	// The Network Condition of the API Token
	Network map[string]interface{} `json:"network,omitempty"`
	// The userId of the user who created the API Token
	UserId string `json:"userId,omitempty"`
	// The client name associated with the API Token
	ClientName string `json:"clientName,omitempty"`
}
