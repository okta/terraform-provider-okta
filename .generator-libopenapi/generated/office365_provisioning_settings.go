// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// Office365ProvisioningSettings represents the Office365ProvisioningSettings schema
// Settings required for the Microsoft Office 365 provisioning connection
type Office365ProvisioningSettings struct {
	// Microsoft Office 365 global administrator password
	AdminPassword string `json:"adminPassword"`
	// Microsoft Office 365 global administrator username
	AdminUsername string `json:"adminUsername"`
}
