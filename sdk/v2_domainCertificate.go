package sdk

type DomainCertificateResource resource

type DomainCertificate struct {
	Certificate      string `json:"certificate,omitempty"`
	CertificateChain string `json:"certificateChain,omitempty"`
	PrivateKey       string `json:"privateKey,omitempty"`
	Type             string `json:"type,omitempty"`
}
