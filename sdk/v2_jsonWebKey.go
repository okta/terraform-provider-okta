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
	X5c         []string    `json:"x5c,omitempty"`
	X5t         string      `json:"x5t,omitempty"`
	X5tS256     string      `json:"x5t#S256,omitempty"`
	X5u         string      `json:"x5u,omitempty"`
}
