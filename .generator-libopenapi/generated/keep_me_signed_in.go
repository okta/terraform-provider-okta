// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// KeepMeSignedIn represents the KeepMeSignedIn schema
// <x-lifecycle-container><x-lifecycle class="oie"></x-lifecycle></x-lifecycle-container>Controls how often the post-authentication prompt is presented to users
type KeepMeSignedIn struct {
	// Whether the post-authentication [Keep Me Signed In (KMSI)](https://help.okta.com/oie/en-us/content/topics/security/stay-signed-in.htm) flow is allowed
	PostAuth string `json:"postAuth,omitempty"`
	// If allowed, how often to display the post-authentication Keep Me Signed In prompt
	PostAuthPromptFrequency interface{} `json:"postAuthPromptFrequency,omitempty"`
}
