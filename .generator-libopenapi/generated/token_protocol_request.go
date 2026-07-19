// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// TokenProtocolRequest represents the TokenProtocolRequest schema
// Details of the token request
type TokenProtocolRequest struct {
	// The authorization response mode
	ResponseMode string `json:"response_mode,omitempty"`
	// The authorization response type
	ResponseType string `json:"response_type,omitempty"`
	// The scopes requested
	Scope string `json:"scope,omitempty"`
	State string `json:"state,omitempty"`
	// The ID of the client associated with the token
	ClientId string `json:"client_id,omitempty"`
	GrantType interface{} `json:"grant_type,omitempty"`
	// Specifies the callback location where the authorization was sent
	RedirectUri string `json:"redirect_uri,omitempty"`
}
