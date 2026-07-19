// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// DevicePostureCheck represents the DevicePostureCheck schema
type DevicePostureCheck struct {
	// Time the device posture check was created
	CreatedDate string `json:"createdDate,omitempty"`
	// Description of the device posture check
	Description string `json:"description,omitempty"`
	// OSQuery for the device posture check
	Query string `json:"query,omitempty"`
	RemediationSettings interface{} `json:"remediationSettings,omitempty"`
	// Unique name of the device posture check
	VariableName string `json:"variableName,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	// User who created the device posture check
	CreatedBy string `json:"createdBy,omitempty"`
	// The ID of the device posture check
	ID string `json:"id,omitempty"`
	// Time the device posture check was updated
	LastUpdate string `json:"lastUpdate,omitempty"`
	// User who updated the device posture check
	LastUpdatedBy string `json:"lastUpdatedBy,omitempty"`
	MappingType interface{} `json:"mappingType,omitempty"`
	// Display name of the device posture check
	Name string `json:"name,omitempty"`
	Platform interface{} `json:"platform,omitempty"`
	Type interface{} `json:"type,omitempty"`
}
