// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UniversalLogoutCapability represents the UniversalLogoutCapability schema
// Universal Logout capability configuration with embedded protocol details.  **Protocol validation rule:** - The `actions` parameter is required when `supportedProtocols` is set to `ACTIONS`.
type UniversalLogoutCapability struct {
	// Configuration for the Actions protocol. This parameter is required when `supportedProtocols` is set to `ACTIONS`.
	Actions []interface{} `json:"actions,omitempty"`
	// Universal Logout capability identifier
	Capability string `json:"capability"`
	// List of supported logout protocols
	SupportedProtocols []string `json:"supportedProtocols"`
}
