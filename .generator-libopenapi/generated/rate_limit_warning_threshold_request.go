// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// RateLimitWarningThresholdRequest represents the RateLimitWarningThresholdRequest schema
type RateLimitWarningThresholdRequest struct {
	// The threshold value (percentage) of a rate limit that, when exceeded, triggers a warning notification. By default, this value is 90 for Workforce orgs and 60 for CIAM orgs.
	WarningThreshold int `json:"warningThreshold"`
}
