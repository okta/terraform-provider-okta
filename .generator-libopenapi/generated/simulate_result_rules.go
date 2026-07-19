// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SimulateResultRules represents the SimulateResultRules schema
type SimulateResultRules struct {
	// List of all conditions involved for this rule evaluation
	Conditions []interface{} `json:"conditions,omitempty"`
	// The unique ID number of the policy rule
	ID string `json:"id,omitempty"`
	// The name of the policy rule
	Name string `json:"name,omitempty"`
	Status interface{} `json:"status,omitempty"`
}
