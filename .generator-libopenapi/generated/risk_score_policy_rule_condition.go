// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// RiskScorePolicyRuleCondition represents the RiskScorePolicyRuleCondition schema
// Specifies a particular level of risk to match on
type RiskScorePolicyRuleCondition struct {
	// The level to match
	Level string `json:"level"`
	// <x-lifecycle-container><x-lifecycle class="oie"></x-lifecycle></x-lifecycle-container>The minimum risk level to match. Only used in a Session Violation Detection (`SESSION_VIOLATION_DETECTION`) pol...
	MinRiskLevel string `json:"minRiskLevel,omitempty"`
}
