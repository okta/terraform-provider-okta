package sdk

import (
	"time"
)

type Csr struct {
	Created *time.Time `json:"created,omitempty"`
	Csr     string     `json:"csr,omitempty"`
	Id      string     `json:"id,omitempty"`
	Kty     string     `json:"kty,omitempty"`
}
