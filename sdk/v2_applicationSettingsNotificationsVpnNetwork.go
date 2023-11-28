// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type ApplicationSettingsNotificationsVpnNetwork struct {
	Connection string   `json:"connection,omitempty"`
	Exclude    []string `json:"exclude,omitempty"`
	Include    []string `json:"include,omitempty"`
}
