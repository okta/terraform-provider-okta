// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type DomainCertificateMetadata struct {
	Expiration  string `json:"expiration,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
	Subject     string `json:"subject,omitempty"`
}
