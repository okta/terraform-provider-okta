// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import (
	"time"
)

type Role struct {
	Embedded       interface{} `json:"_embedded,omitempty"`
	Links          interface{} `json:"_links,omitempty"`
	AssignmentType string      `json:"assignmentType,omitempty"`
	Created        *time.Time  `json:"created,omitempty"`
	Description    string      `json:"description,omitempty"`
	Id             string      `json:"id,omitempty"`
	Label          string      `json:"label,omitempty"`
	LastUpdated    *time.Time  `json:"lastUpdated,omitempty"`
	Status         string      `json:"status,omitempty"`
	Type           string      `json:"type,omitempty"`
}
