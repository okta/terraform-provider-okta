package sdk

type OrgPreferencesResource resource

type OrgPreferences struct {
	Links             interface{} `json:"_links,omitempty"`
	ShowEndUserFooter *bool       `json:"showEndUserFooter,omitempty"`
}
