// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ApplicationSettings represents the ApplicationSettings schema
// App settings
type ApplicationSettings struct {
	Notes interface{} `json:"notes,omitempty"`
	Notifications interface{} `json:"notifications,omitempty"`
	// The entitlement management opt-in status for the app
	EmOptInStatus string `json:"emOptInStatus,omitempty"`
	// Identifies an additional identity store app, if your app supports it. The `identityStoreId` value must be a valid identity store app ID. This identity store app must be created in the same org as y...
	IdentityStoreId string `json:"identityStoreId,omitempty"`
	// Controls whether Okta automatically assigns users to the app based on the user's role or group membership.
	ImplicitAssignment bool `json:"implicitAssignment,omitempty"`
	// Identifier of an inline hook. Inline hooks are outbound calls from Okta to your own custom code, triggered at specific points in Okta process flows. They allow you to integrate custom functionality...
	InlineHookId string `json:"inlineHookId,omitempty"`
}
