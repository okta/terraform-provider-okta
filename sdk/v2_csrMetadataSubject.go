// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type CsrMetadataSubject struct {
	CommonName             string `json:"commonName,omitempty"`
	CountryName            string `json:"countryName,omitempty"`
	LocalityName           string `json:"localityName,omitempty"`
	OrganizationName       string `json:"organizationName,omitempty"`
	OrganizationalUnitName string `json:"organizationalUnitName,omitempty"`
	StateOrProvinceName    string `json:"stateOrProvinceName,omitempty"`
}
