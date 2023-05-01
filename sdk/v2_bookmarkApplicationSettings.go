package sdk

type BookmarkApplicationSettings struct {
	App                *BookmarkApplicationSettingsApplication `json:"app,omitempty"`
	ImplicitAssignment *bool                                   `json:"implicitAssignment,omitempty"`
	InlineHookId       string                                  `json:"inlineHookId,omitempty"`
	Notes              *ApplicationSettingsNotes               `json:"notes,omitempty"`
	Notifications      *ApplicationSettingsNotifications       `json:"notifications,omitempty"`
}
