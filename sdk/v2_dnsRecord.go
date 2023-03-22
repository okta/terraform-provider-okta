package sdk

type DNSRecord struct {
	Expiration string   `json:"expiration,omitempty"`
	Fqdn       string   `json:"fqdn,omitempty"`
	RecordType string   `json:"recordType,omitempty"`
	Values     []string `json:"values,omitempty"`
}
