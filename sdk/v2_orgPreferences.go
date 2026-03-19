// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type OrgPreferencesResource resource

type OrgPreferences struct {
	Links             interface{} `json:"_links,omitempty"`
	ShowEndUserFooter *bool       `json:"showEndUserFooter,omitempty"`
}
