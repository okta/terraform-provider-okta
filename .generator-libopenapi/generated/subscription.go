// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// Subscription represents the Subscription schema
type Subscription struct {
	// An array of sources send notifications to users. > **Note**: Currently, Okta only allows `email` channels.
	Channels []string `json:"channels,omitempty"`
	NotificationType interface{} `json:"notificationType,omitempty"`
	Status interface{} `json:"status,omitempty"`
	// Discoverable resources related to the subscription
	Links map[string]interface{} `json:"_links,omitempty"`
}
