// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PolicyContext represents the PolicyContext schema
type PolicyContext struct {
	Device map[string]interface{} `json:"device,omitempty"`
	// An array of Group IDs for the simulate operation. Only user IDs or Group IDs are allowed, not both.
	Groups map[string]interface{} `json:"groups"`
	// The network rule condition, zone, or IP address
	Ip string `json:"ip,omitempty"`
	// The risk rule condition level
	Risk map[string]interface{} `json:"risk,omitempty"`
	// The user ID for the simulate operation. Only user IDs or Group IDs are allowed, not both.
	User map[string]interface{} `json:"user"`
	// The zone ID under the network rule condition.
	Zones map[string]interface{} `json:"zones,omitempty"`
}
