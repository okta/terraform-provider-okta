// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import (
	"time"
)

type EmailTemplateCustomization struct {
	Links       interface{} `json:"_links,omitempty"`
	Body        string      `json:"body,omitempty"`
	Created     *time.Time  `json:"created,omitempty"`
	Id          string      `json:"id,omitempty"`
	IsDefault   *bool       `json:"isDefault,omitempty"`
	Language    string      `json:"language,omitempty"`
	LastUpdated *time.Time  `json:"lastUpdated,omitempty"`
	Subject     string      `json:"subject,omitempty"`
}
