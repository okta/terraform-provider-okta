// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type DNSRecord struct {
	Expiration string   `json:"expiration,omitempty"`
	Fqdn       string   `json:"fqdn,omitempty"`
	RecordType string   `json:"recordType,omitempty"`
	Values     []string `json:"values,omitempty"`
}
