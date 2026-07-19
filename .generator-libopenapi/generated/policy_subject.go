// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PolicySubject represents the PolicySubject schema
// Specifies the behavior for establishing, validating, and matching a username for an IdP user
type PolicySubject struct {
	// Optional [regular expression pattern](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Regular_expressions) used to filter untrusted IdP usernames. * As a best security practice, y...
	Filter string `json:"filter,omitempty"`
	// Okta user profile attribute for matching a transformed IdP username. Only for matchType `CUSTOM_ATTRIBUTE`. The `matchAttribute` must be a valid Okta user profile attribute of one of the following ...
	MatchAttribute string `json:"matchAttribute,omitempty"`
	MatchType interface{} `json:"matchType,omitempty"`
	UserNameTemplate interface{} `json:"userNameTemplate,omitempty"`
}
