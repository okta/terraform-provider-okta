// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type ApplicationSettings struct {
	App                *ApplicationSettingsApplication   `json:"app,omitempty"`
	ImplicitAssignment *bool                             `json:"implicitAssignment,omitempty"`
	InlineHookId       string                            `json:"inlineHookId,omitempty"`
	Notes              *ApplicationSettingsNotes         `json:"notes,omitempty"`
	Notifications      *ApplicationSettingsNotifications `json:"notifications,omitempty"`
}
