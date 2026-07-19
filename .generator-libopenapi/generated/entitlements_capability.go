// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// EntitlementsCapability represents the EntitlementsCapability schema
// Entitlements capability configuration with embedded protocol details.  **Protocol validation rule:** - The `actions` parameter is required when `supportedProtocols` is set to `ACTIONS`.
type EntitlementsCapability struct {
	// Configuration for the Actions protocol. This parameter is required when `supportedProtocols` is set to `ACTIONS`.
	Actions []interface{} `json:"actions,omitempty"`
	// Entitlements capability identifier
	Capability string `json:"capability"`
	// List of supported entitlements protocols
	SupportedProtocols []string `json:"supportedProtocols"`
}
