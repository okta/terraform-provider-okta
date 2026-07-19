// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ApplicationVisibility represents the ApplicationVisibility schema
// Specifies visibility settings for the app
type ApplicationVisibility struct {
	// Automatically sign in when user lands on the sign-in page
	AutoSubmitToolbar bool `json:"autoSubmitToolbar,omitempty"`
	Hide interface{} `json:"hide,omitempty"`
	// Links or icons that appear on the End-User Dashboard if they're set to `true`.
	AppLinks map[string]interface{} `json:"appLinks,omitempty"`
	// Automatically signs in to the app when user signs into Okta
	AutoLaunch bool `json:"autoLaunch,omitempty"`
}
