// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// DeviceAssurance represents the DeviceAssurance schema
type DeviceAssurance struct {
	ID string `json:"id,omitempty"`
	LastUpdate string `json:"lastUpdate,omitempty"`
	// Display name of the device assurance policy
	Name string `json:"name,omitempty"`
	Platform interface{} `json:"platform,omitempty"`
	CreatedDate string `json:"createdDate,omitempty"`
	DevicePostureChecks interface{} `json:"devicePostureChecks,omitempty"`
	LastUpdatedBy string `json:"lastUpdatedBy,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	CreatedBy string `json:"createdBy,omitempty"`
	// Represents the remediation mode of this device assurance policy when users are denied access due to device noncompliance
	DisplayRemediationMode string `json:"displayRemediationMode,omitempty"`
	GracePeriod interface{} `json:"gracePeriod,omitempty"`
}
