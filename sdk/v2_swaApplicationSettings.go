// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type SwaApplicationSettings struct {
	App                *SwaApplicationSettingsApplication `json:"app,omitempty"`
	ImplicitAssignment *bool                              `json:"implicitAssignment,omitempty"`
	InlineHookId       string                             `json:"inlineHookId,omitempty"`
	Notes              *ApplicationSettingsNotes          `json:"notes,omitempty"`
	Notifications      *ApplicationSettingsNotifications  `json:"notifications,omitempty"`
}

// type of 'App' field is map[string]interface{}, this is the only difference compared to SwaApplicationSettings
type SwaApplicationSettingsWithJSON struct {
	App                *ApplicationSettingsApplication   `json:"app,omitempty"`
	ImplicitAssignment *bool                             `json:"implicitAssignment,omitempty"`
	InlineHookId       string                            `json:"inlineHookId,omitempty"`
	Notes              *ApplicationSettingsNotes         `json:"notes,omitempty"`
	Notifications      *ApplicationSettingsNotifications `json:"notifications,omitempty"`
}
