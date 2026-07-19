// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// Session represents the Session schema
type Session struct {
	// Authentication method reference
	Amr []interface{} `json:"amr,omitempty"`
	// A unique key for the session
	ID string `json:"id,omitempty"`
	// A timestamp when the user last performed the primary or step-up authentication with a password
	LastPasswordVerification *time.Time `json:"lastPasswordVerification,omitempty"`
	// A unique identifier for the user (username)
	Login string `json:"login,omitempty"`
	// Current session status
	Status interface{} `json:"status,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	// A timestamp when the session expires
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	Idp interface{} `json:"idp,omitempty"`
	// A timestamp when the user last performed multifactor authentication
	LastFactorVerification *time.Time `json:"lastFactorVerification,omitempty"`
	// A unique key for the user
	UserId string `json:"userId,omitempty"`
}
