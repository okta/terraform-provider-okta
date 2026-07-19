// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SingleLogout represents the SingleLogout schema
// Determines if the app supports Single Logout (SLO)
type SingleLogout struct {
	// Whether the application supports SLO
	Enabled bool `json:"enabled,omitempty"`
	// The issuer of the Service Provider that generates the SLO request
	Issuer string `json:"issuer,omitempty"`
	// The location where the logout response is sent
	LogoutUrl string `json:"logoutUrl,omitempty"`
}
