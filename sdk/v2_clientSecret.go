// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import (
	"time"
)

type ClientSecret struct {
	Links        interface{} `json:"_links,omitempty"`
	ClientSecret string      `json:"client_secret,omitempty"`
	Created      *time.Time  `json:"created,omitempty"`
	Id           string      `json:"id,omitempty"`
	LastUpdated  *time.Time  `json:"lastUpdated,omitempty"`
	SecretHash   string      `json:"secret_hash,omitempty"`
	Status       string      `json:"status,omitempty"`
}

func NewClientSecret() *ClientSecret {
	return &ClientSecret{}
}

func (a *ClientSecret) IsApplicationInstance() bool {
	return true
}
