// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type DomainCertificateResource resource

type DomainCertificate struct {
	Certificate      string `json:"certificate,omitempty"`
	CertificateChain string `json:"certificateChain,omitempty"`
	PrivateKey       string `json:"privateKey,omitempty"`
	Type             string `json:"type,omitempty"`
}
