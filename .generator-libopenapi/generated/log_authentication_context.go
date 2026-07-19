// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// LogAuthenticationContext represents the LogAuthenticationContext schema
// All authentication relies on validating one or more credentials that prove the authenticity of the actor's identity. Credentials are sometimes provided by the actor, as is the case with passwords, ...
type LogAuthenticationContext struct {
	CredentialType interface{} `json:"credentialType,omitempty"`
	// A proxy for the actor's [session ID](https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html)
	ExternalSessionId string `json:"externalSessionId,omitempty"`
	// The third-party user interface that the actor authenticates through, if any.
	Interface string `json:"interface,omitempty"`
	Issuer interface{} `json:"issuer,omitempty"`
	AuthenticationProvider interface{} `json:"authenticationProvider,omitempty"`
	// The zero-based step number in the authentication pipeline. Currently unused and always set to `0`.
	AuthenticationStep int `json:"authenticationStep,omitempty"`
	CredentialProvider interface{} `json:"credentialProvider,omitempty"`
}
