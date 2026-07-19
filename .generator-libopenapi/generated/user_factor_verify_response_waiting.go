// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// UserFactorVerifyResponseWaiting represents the UserFactorVerifyResponseWaiting schema
type UserFactorVerifyResponseWaiting struct {
	Links interface{} `json:"_links,omitempty"`
	// Timestamp when the verification expires
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	// Optional display message for factor verification
	FactorMessage string `json:"factorMessage,omitempty"`
	FactorResult interface{} `json:"factorResult,omitempty"`
	Profile map[string]interface{} `json:"profile,omitempty"`
	Embedded interface{} `json:"_embedded,omitempty"`
}
