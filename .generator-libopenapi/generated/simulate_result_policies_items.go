// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SimulateResultPoliciesItems represents the SimulateResultPoliciesItems schema
type SimulateResultPoliciesItems struct {
	Rules []interface{} `json:"rules,omitempty"`
	Status interface{} `json:"status,omitempty"`
	// List of all conditions involved for this policy evaluation
	Conditions []interface{} `json:"conditions,omitempty"`
	// ID of the specified policy type
	ID string `json:"id,omitempty"`
	// Policy name
	Name string `json:"name,omitempty"`
}
