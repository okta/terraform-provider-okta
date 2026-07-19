// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// DeviceSignalCollectionPolicyRuleDeviceSignalCollection represents the DeviceSignalCollectionPolicyRuleDeviceSignalCollection schema
// Specifies how device context is collected when a user attempts to sign in
type DeviceSignalCollectionPolicyRuleDeviceSignalCollection struct {
	// Contains the device context provider configuration
	DeviceContextProviders []interface{} `json:"deviceContextProviders,omitempty"`
}
