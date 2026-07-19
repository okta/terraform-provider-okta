// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SimulatePolicyBody represents the SimulatePolicyBody schema
// The request body required for a simulate policy operation
type SimulatePolicyBody struct {
	PolicyContext interface{} `json:"policyContext,omitempty"`
	// Supported policy types for a simulate operation. The default value, `null`, returns all types.
	PolicyTypes []interface{} `json:"policyTypes,omitempty"`
	// The application instance ID for a simulate operation
	AppInstance string `json:"appInstance"`
}
