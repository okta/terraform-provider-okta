package sdk

import (
	"time"
)

type OAuth2RefreshToken struct {
	Embedded    interface{}  `json:"_embedded,omitempty"`
	Links       interface{}  `json:"_links,omitempty"`
	ClientId    string       `json:"clientId,omitempty"`
	Created     *time.Time   `json:"created,omitempty"`
	CreatedBy   *OAuth2Actor `json:"createdBy,omitempty"`
	ExpiresAt   *time.Time   `json:"expiresAt,omitempty"`
	Id          string       `json:"id,omitempty"`
	Issuer      string       `json:"issuer,omitempty"`
	LastUpdated *time.Time   `json:"lastUpdated,omitempty"`
	Scopes      []string     `json:"scopes,omitempty"`
	Status      string       `json:"status,omitempty"`
	UserId      string       `json:"userId,omitempty"`
}
