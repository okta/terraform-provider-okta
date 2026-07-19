// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// DbscStartResponse represents the DbscStartResponse schema
type DbscStartResponse struct {
	Credentials []interface{} `json:"credentials"`
	// URL to call for cookie refresh
	RefreshUrl string `json:"refresh_url"`
	Scope interface{} `json:"scope"`
	// The session identifier for this DBSC binding. Use this value in the `Sec-Secure-Session-Id` header for refresh requests.
	SessionIdentifier string `json:"session_identifier"`
}
