// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// GroupPushMapping represents the GroupPushMapping schema
type GroupPushMapping struct {
	AppConfig map[string]interface{} `json:"appConfig,omitempty"`
	// Timestamp when the group push mapping was created
	Created *time.Time `json:"created,omitempty"`
	// The error message summary if the latest push failed
	ErrorSummary string `json:"errorSummary,omitempty"`
	// Timestamp when the group push mapping was pushed
	LastPush *time.Time `json:"lastPush,omitempty"`
	// The ID of the source group for the group push mapping
	SourceGroupId string `json:"sourceGroupId,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	// The ID of the group push mapping
	ID string `json:"id,omitempty"`
	// Timestamp when the group push mapping was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// The status of the group push mapping
	Status string `json:"status,omitempty"`
	// The ID of the target group for the group push mapping
	TargetGroupId string `json:"targetGroupId,omitempty"`
}
