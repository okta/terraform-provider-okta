// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SimulatePolicyEvaluations represents the SimulatePolicyEvaluations schema
type SimulatePolicyEvaluations struct {
	Status interface{} `json:"status,omitempty"`
	// A list of undefined but not matched policies and rules
	Undefined map[string]interface{} `json:"undefined,omitempty"`
	// A list of evaluated but not matched policies and rules
	Evaluated map[string]interface{} `json:"evaluated,omitempty"`
	// The policy type of the simulate operation
	PolicyType []interface{} `json:"policyType,omitempty"`
	Result interface{} `json:"result,omitempty"`
}
