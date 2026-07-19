// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// UserFactorActivatePush represents the UserFactorActivatePush schema
// Activation requests have a short lifetime and expire if the activation isn't completed before the indicated timestamp. If the activation expires, use the returned `activate` link to restart the pro...
type UserFactorActivatePush struct {
	// Timestamp when the factor verification attempt expires
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	FactorResult interface{} `json:"factorResult,omitempty"`
}
