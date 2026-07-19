// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PostAuthSessionPolicyRuleRunWorkflow represents the PostAuthSessionPolicyRuleRunWorkflow schema
type PostAuthSessionPolicyRuleRunWorkflow struct {
	Action string `json:"action,omitempty"`
	// This action runs a workflow
	Workflow map[string]interface{} `json:"workflow,omitempty"`
}
