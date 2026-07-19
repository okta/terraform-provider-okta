// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// Agent represents the Agent schema
// Agent details
type Agent struct {
	// Unix timestamp in milliseconds when the agent last connected to Okta
	LastConnection int64 `json:"lastConnection,omitempty"`
	// Agent name
	Name string `json:"name,omitempty"`
	OperationalStatus interface{} `json:"operationalStatus,omitempty"`
	// Pool ID
	PoolId string `json:"poolId,omitempty"`
	Type interface{} `json:"type,omitempty"`
	// Status message of the agent
	UpdateMessage string `json:"updateMessage,omitempty"`
	UpdateStatus interface{} `json:"updateStatus,omitempty"`
	// Agent version number
	Version string `json:"version,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	// Unique identifier for the agent that's generated during installation
	ID string `json:"id,omitempty"`
	// Determines if an agent is hidden from the Admin Console
	IsHidden bool `json:"isHidden,omitempty"`
	// Determines if the agent is on the latest generally available version
	IsLatestGAedVersion bool `json:"isLatestGAedVersion,omitempty"`
}
