// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ApplicationSettingsNotificationsVpn represents the ApplicationSettingsNotificationsVpn schema
// Sends customizable messages with conditions to end users when a VPN connection is required
type ApplicationSettingsNotificationsVpn struct {
	Network interface{} `json:"network"`
	// An optional URL to a help page to assist your end users in signing in to your company VPN
	HelpUrl string `json:"helpUrl,omitempty"`
	// A VPN requirement message that's displayed to users
	Message string `json:"message,omitempty"`
}
