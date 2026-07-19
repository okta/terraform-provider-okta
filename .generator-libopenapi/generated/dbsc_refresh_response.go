// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// DbscRefreshResponse represents the DbscRefreshResponse schema
type DbscRefreshResponse struct {
	Credentials []interface{} `json:"credentials"`
	// URL to call for cookie refresh
	RefreshUrl string `json:"refresh_url"`
	Scope interface{} `json:"scope"`
	// The session identifier for this DBSC binding
	SessionIdentifier string `json:"session_identifier"`
}
