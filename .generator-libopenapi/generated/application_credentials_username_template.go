// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ApplicationCredentialsUsernameTemplate represents the ApplicationCredentialsUsernameTemplate schema
// The template used to generate the username when the app is assigned through a group or directly to a user
type ApplicationCredentialsUsernameTemplate struct {
	// Type of mapping expression. Empty string is allowed.
	Type string `json:"type,omitempty"`
	// An optional suffix appended to usernames for `BUILT_IN` mapping expressions
	UserSuffix string `json:"userSuffix,omitempty"`
	// Determines if the username is pushed to the app on updates for CUSTOM `type`
	PushStatus string `json:"pushStatus,omitempty"`
	// Mapping expression used to generate usernames.  The following are supported mapping expressions that are used with the `BUILT_IN` template type:  | Name                            | Template Expres...
	Template string `json:"template,omitempty"`
}
