// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AgentPool represents the AgentPool schema
// An agent pool is a collection of agents that serve a common purpose. An agent pool has a unique ID within an org, and contains a collection of agents disjoint to every other agent pool, meaning tha...
type AgentPool struct {
	Agents []interface{} `json:"agents,omitempty"`
	// Number of agents in the pool that are in a disrupted state
	DisruptedAgents int `json:"disruptedAgents,omitempty"`
	// Agent pool ID
	ID string `json:"id,omitempty"`
	// Number of agents in the pool that are in an inactive state
	InactiveAgents int `json:"inactiveAgents,omitempty"`
	// Agent pool name
	Name string `json:"name,omitempty"`
	OperationalStatus interface{} `json:"operationalStatus,omitempty"`
	Type interface{} `json:"type,omitempty"`
	Links interface{} `json:"_links,omitempty"`
}
