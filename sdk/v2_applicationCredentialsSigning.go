package sdk

import (
	"time"
)

type ApplicationCredentialsSigning struct {
	Kid          string     `json:"kid,omitempty"`
	LastRotated  *time.Time `json:"lastRotated,omitempty"`
	NextRotation *time.Time `json:"nextRotation,omitempty"`
	RotationMode string     `json:"rotationMode,omitempty"`
	Use          string     `json:"use,omitempty"`
}
