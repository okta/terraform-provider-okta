// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type ApplicationSettingsNotificationsVpn struct {
	HelpUrl string                                      `json:"helpUrl,omitempty"`
	Message string                                      `json:"message,omitempty"`
	Network *ApplicationSettingsNotificationsVpnNetwork `json:"network,omitempty"`
}
