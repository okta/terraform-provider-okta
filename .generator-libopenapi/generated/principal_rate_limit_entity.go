// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// PrincipalRateLimitEntity represents the PrincipalRateLimitEntity schema
type PrincipalRateLimitEntity struct {
	// The date and time the principle rate limit entity was created
	CreatedDate *time.Time `json:"createdDate,omitempty"`
	// The default percentage of a given concurrency limit threshold that the owning principal can consume
	DefaultConcurrencyPercentage int `json:"defaultConcurrencyPercentage,omitempty"`
	// The default percentage of a given rate limit threshold that the owning principal can consume
	DefaultPercentage int `json:"defaultPercentage,omitempty"`
	// The unique identifier of the principle rate limit entity
	ID string `json:"id,omitempty"`
	// The date and time the principle rate limit entity was last updated
	LastUpdate *time.Time `json:"lastUpdate,omitempty"`
	// The Okta user ID of the user who last updated the principle rate limit entity
	LastUpdatedBy string `json:"lastUpdatedBy,omitempty"`
	// The unique identifier of the Okta org
	OrgId string `json:"orgId,omitempty"`
	// The unique identifier of the principal. This is the ID of the API token or OAuth 2.0 app.
	PrincipalId string `json:"principalId"`
	PrincipalType interface{} `json:"principalType"`
	// The Okta user ID of the user who created the principle rate limit entity
	CreatedBy string `json:"createdBy,omitempty"`
}
