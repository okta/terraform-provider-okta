// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// BotProtectionConfiguration represents the BotProtectionConfiguration schema
// Bot protection configuration for the org
type BotProtectionConfiguration struct {
	Level interface{} `json:"level"`
	Mode interface{} `json:"mode"`
	// An array of authentication flows that have bot protection enabled
	SupportedFlows []interface{} `json:"supportedFlows,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	EnforcementType interface{} `json:"enforcementType,omitempty"`
}
