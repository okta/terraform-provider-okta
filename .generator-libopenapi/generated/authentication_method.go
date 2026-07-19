// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AuthenticationMethod represents the AuthenticationMethod schema
type AuthenticationMethod struct {
	// Indicates if a user is required to be verified with a verification method.
	UserVerification string `json:"userVerification,omitempty"`
	// Indicates which methods can be used for user verification. `userVerificationMethods` can only be used when `userVerification` is `REQUIRED`. `BIOMETRICS` is currently the only supported method.
	UserVerificationMethods []string `json:"userVerificationMethods,omitempty"`
	// Indicates if any secrets or private keys used during authentication must be hardware protected and not exportable. This property is only set for `POSSESSION` constraints.
	HardwareProtection string `json:"hardwareProtection,omitempty"`
	// An ID that identifies the authenticator
	ID string `json:"id,omitempty"`
	// A label that identifies the authenticator
	Key string `json:"key"`
	// Specifies the method used for the authenticator
	Method string `json:"method"`
	// Indicates if phishing-resistant Factors are required. This property is only set for `POSSESSION` constraints
	PhishingResistant string `json:"phishingResistant,omitempty"`
}
