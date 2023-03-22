package sdk

type ApplicationSettingsNotificationsVpn struct {
	HelpUrl string                                      `json:"helpUrl,omitempty"`
	Message string                                      `json:"message,omitempty"`
	Network *ApplicationSettingsNotificationsVpnNetwork `json:"network,omitempty"`
}
