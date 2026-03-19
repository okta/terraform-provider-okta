// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import (
	"time"
)

type JsonWebKey struct {
	Links       interface{} `json:"_links,omitempty"`
	Alg         string      `json:"alg,omitempty"`
	Created     *time.Time  `json:"created,omitempty"`
	E           string      `json:"e,omitempty"`
	ExpiresAt   *time.Time  `json:"expiresAt,omitempty"`
	KeyOps      []string    `json:"key_ops,omitempty"`
	Kid         string      `json:"kid,omitempty"`
	Kty         string      `json:"kty,omitempty"`
	LastUpdated *time.Time  `json:"lastUpdated,omitempty"`
	N           string      `json:"n,omitempty"`
	Status      string      `json:"status,omitempty"`
	Use         string      `json:"use,omitempty"`
	X           string      `json:"x,omitempty"` // NOTE: EC X parameter is undocumented in docs and oas3 but still valid in the API
	X5c         []string    `json:"x5c,omitempty"`
	X5t         string      `json:"x5t,omitempty"`
	X5tS256     string      `json:"x5t#S256,omitempty"`
	X5u         string      `json:"x5u,omitempty"`
	Y           string      `json:"y,omitempty"` // NOTE: EC Y parameter is undocumented in docs and oas3 but still valid in the API
}
