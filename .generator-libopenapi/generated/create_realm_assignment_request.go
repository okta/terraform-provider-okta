// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CreateRealmAssignmentRequest represents the CreateRealmAssignmentRequest schema
type CreateRealmAssignmentRequest struct {
	Actions interface{} `json:"actions,omitempty"`
	Conditions interface{} `json:"conditions,omitempty"`
	// Name of the realm
	Name string `json:"name,omitempty"`
	// The priority of the realm assignment. The lower the number, the higher the priority. This helps resolve conflicts between realm assignments. > **Note:** When you create realm assignments in bulk, r...
	Priority int `json:"priority,omitempty"`
}
