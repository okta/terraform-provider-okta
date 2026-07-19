// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SsoCapability represents the SsoCapability schema
// SSO capability configuration with embedded protocol details
type SsoCapability struct {
	// SSO capability identifier
	Capability string `json:"capability"`
	// List of supported SSO protocols
	SupportedProtocols []string `json:"supportedProtocols"`
}
