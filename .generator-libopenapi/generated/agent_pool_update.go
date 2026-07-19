// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AgentPoolUpdate represents the AgentPoolUpdate schema
// Various information about agent auto-update configuration
type AgentPoolUpdate struct {
	AgentType interface{} `json:"agentType,omitempty"`
	// ID of the agent pool update
	ID string `json:"id,omitempty"`
	// Indicates if the admin is notified about the update
	NotifyAdmin bool `json:"notifyAdmin,omitempty"`
	// Reason for the update
	Reason string `json:"reason,omitempty"`
	Schedule interface{} `json:"schedule,omitempty"`
	// Indicates if auto-update is enabled for the agent pool
	Enabled bool `json:"enabled,omitempty"`
	// Name of the agent pool update
	Name string `json:"name,omitempty"`
	// Specifies the sort order
	SortOrder int `json:"sortOrder,omitempty"`
	Status interface{} `json:"status,omitempty"`
	// The agent version to update to
	TargetVersion string `json:"targetVersion,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	Agents []interface{} `json:"agents,omitempty"`
}
