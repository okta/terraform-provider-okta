// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// EntityRiskPolicyRuleActionRunWorkflow represents the EntityRiskPolicyRuleActionRunWorkflow schema
type EntityRiskPolicyRuleActionRunWorkflow struct {
	Action string `json:"action,omitempty"`
	// This action runs a workflow
	Workflow map[string]interface{} `json:"workflow,omitempty"`
}
