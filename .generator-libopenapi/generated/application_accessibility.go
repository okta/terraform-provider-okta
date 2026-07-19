// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ApplicationAccessibility represents the ApplicationAccessibility schema
// Specifies access settings for the app
type ApplicationAccessibility struct {
	// Custom error page URL for the app
	ErrorRedirectUrl string `json:"errorRedirectUrl,omitempty"`
	// Custom login page URL for the app > **Note:** The `loginRedirectUrl` property is deprecated in Identity Engine. This property is used with the custom app login feature. Orgs that actively use this ...
	LoginRedirectUrl string `json:"loginRedirectUrl,omitempty"`
	// Represents whether the app can be self-assignable by users
	SelfService bool `json:"selfService,omitempty"`
}
