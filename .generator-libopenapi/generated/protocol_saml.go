// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ProtocolSaml represents the ProtocolSaml schema
// Protocol settings for the [SAML 2.0 Authentication Request Protocol](http://docs.oasis-open.org/security/saml/v2.0/saml-core-2.0-os.pdf)
type ProtocolSaml struct {
	Algorithms interface{} `json:"algorithms,omitempty"`
	Credentials interface{} `json:"credentials,omitempty"`
	Endpoints interface{} `json:"endpoints,omitempty"`
	RelayState interface{} `json:"relayState,omitempty"`
	Settings interface{} `json:"settings,omitempty"`
	// SAML 2.0 protocol
	Type string `json:"type,omitempty"`
}
