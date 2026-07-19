// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// ApplicationCredentialsSigning represents the ApplicationCredentialsSigning schema
// App signing key properties > **Note:** Only apps with SAML_2_0, SAML_1_1, WS_FEDERATION, or OPENID_CONNECT `signOnMode` support the key rotation feature.
type ApplicationCredentialsSigning struct {
	Use interface{} `json:"use,omitempty"`
	// Key identifier used for signing assertions > **Note:** Currently, only the X.509 JWK format is supported for apps with SAML_2_0 `signOnMode`.
	Kid string `json:"kid,omitempty"`
	// Timestamp when the signing key was last rotated
	LastRotated *time.Time `json:"lastRotated,omitempty"`
	// The scheduled time for the next signing key rotation
	NextRotation *time.Time `json:"nextRotation,omitempty"`
	// The mode of key rotation
	RotationMode string `json:"rotationMode,omitempty"`
}
