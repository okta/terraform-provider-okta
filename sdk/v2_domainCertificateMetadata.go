package sdk

type DomainCertificateMetadata struct {
	Expiration  string `json:"expiration,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
	Subject     string `json:"subject,omitempty"`
}
