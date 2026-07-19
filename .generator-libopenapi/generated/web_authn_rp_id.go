// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// WebAuthnRpId represents the WebAuthnRpId schema
// The [RP ID](https://www.w3.org/TR/webauthn/#relying-party-identifier) object for WebAuthn configuration  > **Note:** Changing the RP ID `domain` invalidates all existing passkeys for all users. You...
type WebAuthnRpId struct {
	Domain interface{} `json:"domain,omitempty"`
	// Indicates whether the RP ID is active and is used for WebAuthn operations. It can only be set to `true` once the `validationStatus` of the `domain` object is `VERIFIED`. `enabled` can only be `true...
	Enabled bool `json:"enabled,omitempty"`
}
