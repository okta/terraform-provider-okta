// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AgentPoolUpdateSettingResponse represents the AgentPoolUpdateSettingResponse schema
// Settings for auto-update
type AgentPoolUpdateSettingResponse struct {
	// Latest version of the agent
	LatestVersion string `json:"latestVersion,omitempty"`
	// Minimal version of the agent
	MinimalSupportedVersion string `json:"minimalSupportedVersion,omitempty"`
	// ID of the agent pool that the settings apply to
	PoolId string `json:"poolId,omitempty"`
	// Pool name
	PoolName string `json:"poolName,omitempty"`
	ReleaseChannel interface{} `json:"releaseChannel,omitempty"`
	AgentType interface{} `json:"agentType,omitempty"`
	// Continues the update even if some agents fail to update
	ContinueOnError bool `json:"continueOnError,omitempty"`
}
