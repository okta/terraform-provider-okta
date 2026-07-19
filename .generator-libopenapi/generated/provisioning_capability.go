// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ProvisioningCapability represents the ProvisioningCapability schema
// Provisioning capability configuration with embedded protocol details.  *Protocol validation rule:** - The `actions` parameter is required when `supportedProtocols` is set to `ACTIONS`.
type ProvisioningCapability struct {
	// Configuration for the Actions protocol. This parameter is required when `supportedProtocols` is set to `ACTIONS`.
	Actions []interface{} `json:"actions,omitempty"`
	// Provisioning capability identifier
	Capability string `json:"capability"`
	// List of supported provisioning protocols
	SupportedProtocols []string `json:"supportedProtocols"`
}
