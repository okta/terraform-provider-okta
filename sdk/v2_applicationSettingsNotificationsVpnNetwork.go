package sdk

type ApplicationSettingsNotificationsVpnNetwork struct {
	Connection string   `json:"connection,omitempty"`
	Exclude    []string `json:"exclude,omitempty"`
	Include    []string `json:"include,omitempty"`
}
