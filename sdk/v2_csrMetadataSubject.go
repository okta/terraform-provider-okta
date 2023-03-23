package sdk

type CsrMetadataSubject struct {
	CommonName             string `json:"commonName,omitempty"`
	CountryName            string `json:"countryName,omitempty"`
	LocalityName           string `json:"localityName,omitempty"`
	OrganizationName       string `json:"organizationName,omitempty"`
	OrganizationalUnitName string `json:"organizationalUnitName,omitempty"`
	StateOrProvinceName    string `json:"stateOrProvinceName,omitempty"`
}
