// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// MtlsTrustCredentials represents the MtlsTrustCredentials schema
type MtlsTrustCredentials struct {
	// Not used
	Audience string `json:"audience,omitempty"`
	// Description of the certificate issuer
	Issuer string `json:"issuer,omitempty"`
	Kid interface{} `json:"kid,omitempty"`
	Revocation interface{} `json:"revocation,omitempty"`
	// Time in minutes to cache the certificate revocation information  > **Note:** This property isn't supported. Okta now handles CRL caching automatically. As of October 8, 2025, in Preview orgs, and O...
	RevocationCacheLifetime float64 `json:"revocationCacheLifetime,omitempty"`
}
