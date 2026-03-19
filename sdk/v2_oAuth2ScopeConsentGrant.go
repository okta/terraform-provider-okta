// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import (
	"time"
)

type OAuth2ScopeConsentGrant struct {
	Embedded    interface{}  `json:"_embedded,omitempty"`
	Links       interface{}  `json:"_links,omitempty"`
	ClientId    string       `json:"clientId,omitempty"`
	Created     *time.Time   `json:"created,omitempty"`
	CreatedBy   *OAuth2Actor `json:"createdBy,omitempty"`
	Id          string       `json:"id,omitempty"`
	Issuer      string       `json:"issuer,omitempty"`
	LastUpdated *time.Time   `json:"lastUpdated,omitempty"`
	ScopeId     string       `json:"scopeId,omitempty"`
	Source      string       `json:"source,omitempty"`
	Status      string       `json:"status,omitempty"`
	UserId      string       `json:"userId,omitempty"`
}
