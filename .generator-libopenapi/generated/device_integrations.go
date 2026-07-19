// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// DeviceIntegrations represents the DeviceIntegrations schema
type DeviceIntegrations struct {
	Status interface{} `json:"status,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	// The display name of the device integration
	DisplayName string `json:"displayName,omitempty"`
	// The ID of the device integration
	ID string `json:"id,omitempty"`
	Metadata interface{} `json:"metadata,omitempty"`
	Name interface{} `json:"name,omitempty"`
	Platform interface{} `json:"platform,omitempty"`
}
