// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// DNSRecordDomains represents the DNSRecordDomains schema
// DNS TXT and CNAME records to be registered for the Domain
type DNSRecordDomains struct {
	// DNS record name
	Fqdn string `json:"fqdn,omitempty"`
	RecordType interface{} `json:"recordType,omitempty"`
	// DNS record value
	Values []string `json:"values,omitempty"`
	// DNS TXT record expiration
	Expiration string `json:"expiration,omitempty"`
}
