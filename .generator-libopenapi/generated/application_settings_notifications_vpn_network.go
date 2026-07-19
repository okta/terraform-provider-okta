// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ApplicationSettingsNotificationsVpnNetwork represents the ApplicationSettingsNotificationsVpnNetwork schema
// Defines network zones for VPN notification
type ApplicationSettingsNotificationsVpnNetwork struct {
	// Specifies the VPN connection details required to access the app
	Connection string `json:"connection,omitempty"`
	// Defines the IP addresses or network ranges that are excluded from the VPN requirement
	Exclude []string `json:"exclude,omitempty"`
	// Defines the IP addresses or network ranges that are required to use the VPN
	Include []string `json:"include,omitempty"`
}
