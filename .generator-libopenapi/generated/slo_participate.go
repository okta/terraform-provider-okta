// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SloParticipate represents the SloParticipate schema
// Determines if the app participates in Single Logout (SLO)
type SloParticipate struct {
	// Request binding type
	BindingType string `json:"bindingType,omitempty"`
	// Indicates whether the app is allowed to participate in front-channel SLO
	Enabled bool `json:"enabled,omitempty"`
	// URL where Okta sends the logout request
	LogoutRequestUrl string `json:"logoutRequestUrl,omitempty"`
	// Determines whether Okta sends the `SessionIndex` elements in the logout request
	SessionIndexRequired bool `json:"sessionIndexRequired,omitempty"`
}
