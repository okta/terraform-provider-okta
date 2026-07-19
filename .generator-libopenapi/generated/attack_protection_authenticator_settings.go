// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AttackProtectionAuthenticatorSettings represents the AttackProtectionAuthenticatorSettings schema
type AttackProtectionAuthenticatorSettings struct {
	// If true, requires users to verify a possession factor before verifying a knowledge factor when the assurance requires two-factor authentication (2FA).
	VerifyKnowledgeSecondWhen2faRequired bool `json:"verifyKnowledgeSecondWhen2faRequired,omitempty"`
}
