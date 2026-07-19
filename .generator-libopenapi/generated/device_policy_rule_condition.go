// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// DevicePolicyRuleCondition represents the DevicePolicyRuleCondition schema
type DevicePolicyRuleCondition struct {
	Migrated bool `json:"migrated,omitempty"`
	Platform interface{} `json:"platform,omitempty"`
	Rooted bool `json:"rooted,omitempty"`
	TrustLevel interface{} `json:"trustLevel,omitempty"`
}
