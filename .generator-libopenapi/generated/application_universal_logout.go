// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ApplicationUniversalLogout represents the ApplicationUniversalLogout schema
// <div class="x-lifecycle-container"><x-lifecycle class="oie"></x-lifecycle></div> Universal Logout properties for the app. These properties are only returned and can't be updated.
type ApplicationUniversalLogout struct {
	// Indicates whether the app uses a shared identity stack that may cause the user to sign out of other apps by the same company
	IdentityStack string `json:"identityStack,omitempty"`
	// The protocol used for Universal Logout
	Protocol string `json:"protocol,omitempty"`
	// Universal Logout status for the app instance
	Status string `json:"status,omitempty"`
	// Indicates whether the app supports full or partial Universal Logout (UL).
	SupportType string `json:"supportType,omitempty"`
}
