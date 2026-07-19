// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// TokenPayLoad represents the TokenPayLoad schema
type TokenPayLoad struct {
	// The URL of the token inline hook
	Source string `json:"source,omitempty"`
	Data map[string]interface{} `json:"data,omitempty"`
	// The type of inline hook. The token inline hook type is `com.okta.oauth2.tokens.transform`.
	EventType string `json:"eventType,omitempty"`
}
