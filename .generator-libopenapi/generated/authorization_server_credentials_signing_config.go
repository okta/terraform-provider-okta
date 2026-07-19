// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// AuthorizationServerCredentialsSigningConfig represents the AuthorizationServerCredentialsSigningConfig schema
type AuthorizationServerCredentialsSigningConfig struct {
	// The ID of the JSON Web Key used for signing tokens issued by the authorization server
	Kid string `json:"kid,omitempty"`
	// The timestamp when the authorization server started using the `kid` for signing tokens
	LastRotated *time.Time `json:"lastRotated,omitempty"`
	// The timestamp when the authorization server changes the Key for signing tokens. This is only returned when `rotationMode` is set to `AUTO`.
	NextRotation *time.Time `json:"nextRotation,omitempty"`
	RotationMode interface{} `json:"rotationMode,omitempty"`
	Use interface{} `json:"use,omitempty"`
}
