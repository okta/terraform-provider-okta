// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CapabilitiesImportRulesUserCreateAndMatchObject represents the CapabilitiesImportRulesUserCreateAndMatchObject schema
// Rules for matching and creating users
type CapabilitiesImportRulesUserCreateAndMatchObject struct {
	// Determines the attribute to match users
	ExactMatchCriteria string `json:"exactMatchCriteria,omitempty"`
	// Allows user import upon partial matching. Partial matching occurs when the first and last names of an imported user match those of an existing Okta user, even if the username or email attributes do...
	AllowPartialMatch bool `json:"allowPartialMatch,omitempty"`
	// If set to `true`, imported new users are automatically activated.
	AutoActivateNewUsers bool `json:"autoActivateNewUsers,omitempty"`
	// If set to `true`, exact-matched users are automatically confirmed on activation. If set to `false`, exact-matched users need to be confirmed manually.
	AutoConfirmExactMatch bool `json:"autoConfirmExactMatch,omitempty"`
	// If set to `true`, imported new users are automatically confirmed on activation. This doesn't apply to imported users that already exist in Okta.
	AutoConfirmNewUsers bool `json:"autoConfirmNewUsers,omitempty"`
	// If set to `true`, partially matched users are automatically confirmed on activation. If set to `false`, partially matched users need to be confirmed manually.
	AutoConfirmPartialMatch bool `json:"autoConfirmPartialMatch,omitempty"`
}
