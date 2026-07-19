// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ApiServiceCapability represents the ApiServiceCapability schema
// API Service capability configuration with embedded protocol details
type ApiServiceCapability struct {
	// List of supported protocols
	SupportedProtocols []string `json:"supportedProtocols"`
	// API Service capability identifier
	Capability string `json:"capability"`
}
